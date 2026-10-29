// Смоук-тест запущенного сервера: регистрация, синк двух устройств, LWW, блокеры, теги, правила
// целостности при удалении (сервер ничего не удаляет каскадом — сначала дети, потом родитель). Нужна открытая регистрация (AUTH_SIGN_UP_ENABLED=true).
//
//	go run ./cmd/smoke -addr localhost:8081
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	pb "github.com/jedi-tasks/go-server-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var failed bool

func check(name string, ok bool, details ...any) {
	if ok {
		fmt.Println("OK  ", name)
	} else {
		failed = true
		fmt.Println("FAIL", name, details)
	}
}

func must[T interface{ GetError() *pb.Error }](res T, err error) T {
	name := "call"
	if err != nil {
		check(name, false, err)
		os.Exit(1)
	}
	if e := res.GetError(); e != nil {
		check(name, false, e.Message, e.SystemMessage)
		os.Exit(1)
	}
	return res
}

// rejected проверяет, что мутация отклонена с ожидаемым текстом для пользователя
func rejected[T interface{ GetError() *pb.Error }](name string, res T, err error, humanText string) {
	check(name, err == nil && res.GetError() != nil && res.GetError().Message == humanText, err, res.GetError())
}

func id() []byte                            { u := uuid.New(); return u[:] }
func ts(t time.Time) *timestamppb.Timestamp { return timestamppb.New(t) }
func ptr[T any](v T) *T                     { return &v }

func main() {
	addr := flag.String("addr", "localhost:8081", "адрес gRPC-сервера")
	flag.Parse()

	ctx := context.Background()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	auth := pb.NewAuthEndpointClient(conn)
	syncC := pb.NewSyncEndpointClient(conn)
	lists := pb.NewTaskListEndpointClient(conn)
	tasks := pb.NewTaskEndpointClient(conn)
	tags := pb.NewTagEndpointClient(conn)
	blockers := pb.NewBlockerEndpointClient(conn)
	settings := pb.NewSettingsEndpointClient(conn)
	sections := pb.NewListSectionEndpointClient(conn)

	email := fmt.Sprintf("smoke-%d@test.local", time.Now().UnixNano())
	dev := &pb.DeviceInformation{DeviceName: "smoke", ModelName: "Mac", NameOS: pb.OSType_OSX, VersionOS: "15", IpAddress: "127.0.0.1"}
	app := &pb.ApplicationInformation{Version: "1", Build: "1", BundleID: "smoke"}
	up := must(auth.SignUp(ctx, &pb.SignUpRequest{Name: "Smoke", Email: email, Password: "pass-1234", Application: app, Device: dev, DeviceID: "dev-A"}))
	tokA := up.Token.AccessToken
	in := must(auth.SignIn(ctx, &pb.SignInRequest{Email: email, Password: "pass-1234", Application: app, Device: dev, DeviceID: "dev-B"}))
	tokB := in.Token.AccessToken

	doSync := func(tok string, since uint32) *pb.SyncResponse {
		r := must(syncC.Sync(ctx, &pb.SyncRequest{AccessToken: tok, SinceID: since}))
		if r.HasChanges {
			must(syncC.ConfirmSync(ctx, &pb.ConfirmSyncRequest{AccessToken: tok, PendingSyncToken: r.PendingSyncToken}))
		}
		return r
	}

	// Первичный синк: Inbox, 3 приоритета, 7 длительностей
	sA := doSync(tokA, 0)
	check("initial sync: 1 inbox", len(sA.ChangedTaskLists) == 1 && sA.ChangedTaskLists[0].IsInbox, len(sA.ChangedTaskLists))
	check("initial sync: 3 priorities", len(sA.ChangedPriorities) == 3, len(sA.ChangedPriorities))
	check("initial sync: 7 durations", len(sA.ChangedDurationOptions) == 7, len(sA.ChangedDurationOptions))
	cpA := sA.PendingCheckpoint
	sB := doSync(tokB, 0)
	cpB := sB.PendingCheckpoint
	inboxID := sA.ChangedTaskLists[0].Id

	// Задача без списка -> Inbox
	now := time.Now()
	t1 := id()
	must(tasks.CreateTask(ctx, &pb.CreateTaskRequest{AccessToken: tokA, Id: t1, Title: "Купить шины", PlannedDay: ptr("2026-10-30"), DurationMinutes: ptr(int32(30)), CreatedAt: ts(now)}))
	// повтор из очереди — без ошибки и без дубля
	must(tasks.CreateTask(ctx, &pb.CreateTaskRequest{AccessToken: tokA, Id: t1, Title: "Купить шины", CreatedAt: ts(now)}))
	got := must(tasks.GetTasks(ctx, &pb.GetTasksRequest{AccessToken: tokA, Ids: [][]byte{t1}}))
	check("task in inbox", len(got.Tasks) == 1 && string(got.Tasks[0].ListID) == string(inboxID))

	// Устройство B не синхронизировано — мутация отклоняется NeedToSync
	r, err := tasks.UpdateTask(ctx, &pb.UpdateTaskRequest{AccessToken: tokB, Id: t1, Title: ptr("B"), UpdatedAt: ts(now.Add(time.Second))})
	check("sync gate rejects stale device", err == nil && r.Error != nil && r.Error.Category == pb.ErrorCategory_NeedToSync, r.GetError())

	// B синкается и видит задачу
	sB = doSync(tokB, cpB)
	cpB = sB.PendingCheckpoint
	check("device B gets task", len(sB.ChangedTasks) == 1)

	// LWW: B правит в 10:00+2s, A потом шлет более старую правку — игнорируется
	must(tasks.UpdateTask(ctx, &pb.UpdateTaskRequest{AccessToken: tokB, Id: t1, Title: ptr("Новое"), UpdatedAt: ts(now.Add(2 * time.Second))}))
	doSync(tokA, cpA)
	must(tasks.UpdateTask(ctx, &pb.UpdateTaskRequest{AccessToken: tokA, Id: t1, Title: ptr("Старое"), UpdatedAt: ts(now.Add(1 * time.Second))}))
	got = must(tasks.GetTasks(ctx, &pb.GetTasksRequest{AccessToken: tokA, Ids: [][]byte{t1}}))
	check("LWW keeps newer title", got.Tasks[0].Title == "Новое", got.Tasks[0].Title)

	// Сброс поля пустым значением
	must(tasks.UpdateTask(ctx, &pb.UpdateTaskRequest{AccessToken: tokA, Id: t1, PlannedDay: ptr(""), DurationMinutes: ptr(int32(0)), UpdatedAt: ts(now.Add(3 * time.Second))}))
	got = must(tasks.GetTasks(ctx, &pb.GetTasksRequest{AccessToken: tokA, Ids: [][]byte{t1}}))
	check("plannedDay and duration cleared", got.Tasks[0].PlannedDay == nil && got.Tasks[0].DurationMinutes == nil)

	// Блокеры: t2 блокирует t1, цикл t1 -> t2 запрещен. Выполнение t2 блокер не снимает — клиент снимает сам
	t2 := id()
	must(tasks.CreateTask(ctx, &pb.CreateTaskRequest{AccessToken: tokA, Id: t2, Title: "Заказать", CreatedAt: ts(now)}))
	b1 := id()
	must(blockers.CreateBlocker(ctx, &pb.CreateBlockerRequest{AccessToken: tokA, Id: b1, TaskID: t1, BlockingTaskID: t2, SinceDay: "2026-10-29", UpdatedAt: ts(now)}))
	cr, err := blockers.CreateBlocker(ctx, &pb.CreateBlockerRequest{AccessToken: tokA, Id: id(), TaskID: t2, BlockingTaskID: t1, SinceDay: "2026-10-29", UpdatedAt: ts(now)})
	check("cycle rejected", err == nil && cr.Error != nil, cr.GetError())
	must(tasks.SetTaskCompleted(ctx, &pb.SetTaskCompletedRequest{AccessToken: tokA, Id: t2, Completed: true, UpdatedAt: ts(now.Add(4 * time.Second))}))
	bl := must(blockers.GetBlockers(ctx, &pb.GetBlockersRequest{AccessToken: tokA, TaskIDs: [][]byte{t1}}))
	check("blocker not resolved by server", len(bl.Blockers) == 1 && bl.Blockers[0].ResolvedAt == nil)
	must(blockers.ResolveBlocker(ctx, &pb.ResolveBlockerRequest{AccessToken: tokA, Id: b1, UpdatedAt: ts(now.Add(4500 * time.Millisecond))}))
	bl = must(blockers.GetBlockers(ctx, &pb.GetBlockersRequest{AccessToken: tokA, TaskIDs: [][]byte{t1}}))
	check("blocker resolved by client", len(bl.Blockers) == 1 && bl.Blockers[0].ResolvedAt != nil)
	dt, err := tasks.DeleteTask(ctx, &pb.DeleteTaskRequest{AccessToken: tokA, Id: t2, UpdatedAt: ts(now.Add(time.Hour))})
	rejected("task with blocker can't be deleted", dt, err, "Задача участвует в блокерах — сначала удалите их")

	// Теги
	tg := id()
	must(tags.CreateTag(ctx, &pb.CreateTagRequest{AccessToken: tokA, Id: tg, Name: "велосипед", UpdatedAt: ts(now)}))
	must(tags.AddTagToTask(ctx, &pb.AddTagToTaskRequest{AccessToken: tokA, TaskID: t1, TagID: tg, UpdatedAt: ts(now.Add(5 * time.Second))}))
	dtg, err := tags.DeleteTag(ctx, &pb.DeleteTagRequest{AccessToken: tokA, Id: tg, UpdatedAt: ts(now.Add(time.Hour))})
	rejected("linked tag can't be deleted", dtg, err, "Тег привязан к задачам — сначала отвяжите его")

	// Раздел и удаление списка: родитель с живыми детьми не удаляется, сначала удаляем детей
	l := id()
	must(lists.CreateTaskList(ctx, &pb.CreateTaskListRequest{AccessToken: tokA, Id: l, Name: "Работа", ColorHex: "#E5484D", Icon: "hammer.fill", SortOrder: 1, SortMode: pb.SortMode_Manual, CreatedAt: ts(now)}))
	sec := id()
	must(sections.CreateListSection(ctx, &pb.CreateListSectionRequest{AccessToken: tokA, Id: sec, ListID: l, Name: "Срочно", UpdatedAt: ts(now)}))
	t3 := id()
	must(tasks.CreateTask(ctx, &pb.CreateTaskRequest{AccessToken: tokA, Id: t3, ListID: l, SectionID: sec, Title: "Отчет", CreatedAt: ts(now)}))
	dl, err := lists.DeleteTaskList(ctx, &pb.DeleteTaskListRequest{AccessToken: tokA, Id: inboxID, UpdatedAt: ts(now.Add(time.Hour))})
	rejected("inbox can't be deleted", dl, err, "Inbox удалить нельзя")
	dl, err = lists.DeleteTaskList(ctx, &pb.DeleteTaskListRequest{AccessToken: tokA, Id: l, UpdatedAt: ts(now.Add(time.Hour))})
	rejected("list with tasks can't be deleted", dl, err, "В списке есть задачи — сначала удалите их или перенесите в другой список")
	ds, err := sections.DeleteListSection(ctx, &pb.DeleteListSectionRequest{AccessToken: tokA, Id: sec, UpdatedAt: ts(now.Add(time.Hour))})
	rejected("section with tasks can't be deleted", ds, err, "В разделе есть задачи — сначала удалите их или перенесите в другой раздел")
	must(tasks.DeleteTask(ctx, &pb.DeleteTaskRequest{AccessToken: tokA, Id: t3, UpdatedAt: ts(now.Add(6 * time.Second))}))
	dl, err = lists.DeleteTaskList(ctx, &pb.DeleteTaskListRequest{AccessToken: tokA, Id: l, UpdatedAt: ts(now.Add(time.Hour))})
	rejected("list with sections can't be deleted", dl, err, "В списке есть разделы — сначала удалите их")
	must(sections.DeleteListSection(ctx, &pb.DeleteListSectionRequest{AccessToken: tokA, Id: sec, UpdatedAt: ts(now.Add(7 * time.Second))}))
	must(lists.DeleteTaskList(ctx, &pb.DeleteTaskListRequest{AccessToken: tokA, Id: l, UpdatedAt: ts(now.Add(8 * time.Second))}))
	ct, err := tasks.CreateTask(ctx, &pb.CreateTaskRequest{AccessToken: tokA, Id: id(), ListID: l, Title: "В удаленный", CreatedAt: ts(now.Add(9 * time.Second))})
	rejected("task in deleted list rejected", ct, err, "Список удалён — выберите другой список")

	// Настройки
	must(settings.UpdateSettings(ctx, &pb.UpdateSettingsRequest{AccessToken: tokA, DailyCapacityMinutes: ptr(uint32(300)), UpdatedAt: ts(now)}))

	// Устройство B получает всё одним синком
	sB = doSync(tokB, cpB)
	check("B: task list deleted", len(sB.DeletedTaskListIDs) == 1, len(sB.DeletedTaskListIDs))
	check("B: section deleted", len(sB.DeletedListSectionIDs) == 1, len(sB.DeletedListSectionIDs))
	check("B: task deleted before its list", len(sB.DeletedTaskIDs) == 1, len(sB.DeletedTaskIDs))
	check("B: tag link added", len(sB.AddedTasksToTags) == 1, len(sB.AddedTasksToTags))
	check("B: blocker resolved", len(sB.ChangedBlockers) == 1 && sB.ChangedBlockers[0].ResolvedAt != nil)
	check("B: settings changed", sB.ChangedSettings != nil && sB.ChangedSettings.DailyCapacityMinutes == 300)
	check("B: tasks t1,t2 changed", len(sB.ChangedTasks) == 2, len(sB.ChangedTasks))

	// Пустой повторный синк
	empty := doSync(tokB, sB.PendingCheckpoint)
	check("B: nothing more", !empty.HasChanges)

	if failed {
		os.Exit(1)
	}
}
