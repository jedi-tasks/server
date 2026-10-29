package integrationTest

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	taskModel "server/internal/modules/task/model"
	taskListModel "server/internal/modules/taskList/model"
)

// Тексты ошибок целостности, которые видит пользователь
const (
	textListHasTasks       = "В списке есть задачи — сначала удалите их или перенесите в другой список"
	textListHasSections    = "В списке есть разделы — сначала удалите их"
	textSectionHasTasks    = "В разделе есть задачи — сначала удалите их или перенесите в другой раздел"
	textPriorityHasTasks   = "У приоритета есть задачи — сначала удалите их или смените им приоритет"
	textTagHasLinks        = "Тег привязан к задачам — сначала отвяжите его"
	textTaskHasTags        = "К задаче привязаны теги — сначала отвяжите их"
	textTaskHasBlockers    = "Задача участвует в блокерах — сначала удалите их"
	textListDeleted        = "Список удалён — выберите другой список"
	textSectionDeleted     = "Раздел удалён — выберите другой раздел"
	textPriorityDeleted    = "Приоритет удалён — выберите другой приоритет"
	textTaskDeleted        = "Задача удалена"
	textTagDeleted         = "Тег удалён"
	textInboxCantBeDeleted = "Inbox удалить нельзя"
)

// TestDeleteListWithTasks — список с живой задачей не удаляется, после удаления задачи удаляется
func TestDeleteListWithTasks(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	taskID := f.createTask(listID, nil, nil)

	f.rejected(f.deleteList(listID), textListHasTasks)
	f.ok(f.deleteTask(taskID))
	f.ok(f.deleteList(listID))
}

// TestDeleteListWithSections — список с живым разделом не удаляется, после удаления раздела удаляется
func TestDeleteListWithSections(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	sectionID := f.createSection(listID)

	f.rejected(f.deleteList(listID), textListHasSections)
	f.ok(f.deleteSection(sectionID))
	f.ok(f.deleteList(listID))
}

// TestDeleteInbox — Inbox удалить нельзя
func TestDeleteInbox(t *testing.T) {
	f := newFixture(t)
	f.rejected(f.deleteList(f.inboxID), textInboxCantBeDeleted)
}

// TestDeleteSectionWithTasks — раздел с живой задачей не удаляется; после переноса задачи удаляется
func TestDeleteSectionWithTasks(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	sectionID := f.createSection(listID)
	taskID := f.createTask(listID, &sectionID, nil)

	f.rejected(f.deleteSection(sectionID), textSectionHasTasks)

	// Убираем задачу из раздела — ссылка сброшена, раздел удаляется
	noSection := uuid.Nil
	f.ok(f.task.UpdateTask(f.ctx, taskModel.UpdateTaskReq{Necessary: f.necessary, ID: taskID, SectionID: &noSection, UpdatedAt: f.next()})) //nolint:exhaustruct
	f.ok(f.deleteSection(sectionID))
}

// TestDeletePriorityWithTasks — приоритет живой задачи не удаляется, после удаления задачи удаляется
func TestDeletePriorityWithTasks(t *testing.T) {
	f := newFixture(t)
	priorityID := f.createPriority()
	taskID := f.createTask(f.inboxID, nil, &priorityID)

	f.rejected(f.deletePriority(priorityID), textPriorityHasTasks)
	f.ok(f.deleteTask(taskID))
	f.ok(f.deletePriority(priorityID))
}

// TestDeleteTagWithLinks — привязанный тег не удаляется, после отвязки удаляется
func TestDeleteTagWithLinks(t *testing.T) {
	f := newFixture(t)
	taskID := f.createTask(f.inboxID, nil, nil)
	tagID := f.createTag()
	f.ok(f.linkTag(taskID, tagID))

	f.rejected(f.deleteTag(tagID), textTagHasLinks)
	f.ok(f.unlinkTag(taskID, tagID))
	f.ok(f.deleteTag(tagID))
}

// TestDeleteTaskWithTags — задача с привязанным тегом не удаляется, после отвязки удаляется
func TestDeleteTaskWithTags(t *testing.T) {
	f := newFixture(t)
	taskID := f.createTask(f.inboxID, nil, nil)
	tagID := f.createTag()
	f.ok(f.linkTag(taskID, tagID))

	f.rejected(f.deleteTask(taskID), textTaskHasTags)
	f.ok(f.unlinkTag(taskID, tagID))
	f.ok(f.deleteTask(taskID))
}

// TestDeleteTaskWithBlockers — задача не удаляется, пока она заблокирована или блокирует другую
func TestDeleteTaskWithBlockers(t *testing.T) {
	f := newFixture(t)
	blockedID := f.createTask(f.inboxID, nil, nil)
	blockingID := f.createTask(f.inboxID, nil, nil)
	blockerID := uuid.New()
	f.ok(f.createBlocker(blockerID, blockedID, blockingID))

	f.rejected(f.deleteTask(blockedID), textTaskHasBlockers)
	f.rejected(f.deleteTask(blockingID), textTaskHasBlockers)

	f.ok(f.deleteBlocker(blockerID))
	f.ok(f.deleteTask(blockedID))
	f.ok(f.deleteTask(blockingID))
}

// TestCompletedTaskKeepsBlockers — выполнение задачи не снимает ее блокеры: это делает клиент
func TestCompletedTaskKeepsBlockers(t *testing.T) {
	f := newFixture(t)
	blockedID := f.createTask(f.inboxID, nil, nil)
	blockingID := f.createTask(f.inboxID, nil, nil)
	blockerID := uuid.New()
	f.ok(f.createBlocker(blockerID, blockedID, blockingID))

	f.ok(f.task.SetTaskCompleted(f.ctx, taskModel.SetTaskCompletedReq{Necessary: f.necessary, ID: blockingID, Completed: true, UpdatedAt: f.next()}))

	unresolved := f.count(sq.Select("COUNT(*)").From("tasks.blockers").Where(sq.Eq{"id": blockerID, "resolved_at": nil}))
	if unresolved != 1 {
		t.Fatalf("блокер выполненной задачи снят сервером")
	}
}

// TestReferenceToDeletedList — задачу и раздел нельзя создать в удаленном списке или перенести в него
func TestReferenceToDeletedList(t *testing.T) {
	f := newFixture(t)
	deletedListID := f.createList()
	f.ok(f.deleteList(deletedListID))

	f.rejected(f.task.CreateTask(f.ctx, f.taskReq(uuid.New(), deletedListID, nil, nil)), textListDeleted)
	f.rejected(f.listSection.CreateListSection(f.ctx, f.sectionReq(uuid.New(), deletedListID)), textListDeleted)

	taskID := f.createTask(f.inboxID, nil, nil)
	f.rejected(f.task.UpdateTask(f.ctx, taskModel.UpdateTaskReq{Necessary: f.necessary, ID: taskID, ListID: &deletedListID, UpdatedAt: f.next()}), textListDeleted) //nolint:exhaustruct

	// Отклоненные мутации ничего не записали
	if n := f.count(sq.Select("COUNT(*)").From("tasks.tasks").Where(sq.Eq{"list_id": deletedListID})); n != 0 {
		t.Fatalf("в удаленном списке появились задачи: %d", n)
	}
}

// TestReferenceToDeletedSectionAndPriority — задаче нельзя назначить удаленный раздел или приоритет
func TestReferenceToDeletedSectionAndPriority(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	sectionID := f.createSection(listID)
	priorityID := f.createPriority()
	f.ok(f.deleteSection(sectionID))
	f.ok(f.deletePriority(priorityID))

	f.rejected(f.task.CreateTask(f.ctx, f.taskReq(uuid.New(), listID, &sectionID, nil)), textSectionDeleted)
	f.rejected(f.task.CreateTask(f.ctx, f.taskReq(uuid.New(), listID, nil, &priorityID)), textPriorityDeleted)

	taskID := f.createTask(listID, nil, nil)
	f.rejected(f.task.UpdateTask(f.ctx, taskModel.UpdateTaskReq{Necessary: f.necessary, ID: taskID, SectionID: &sectionID, UpdatedAt: f.next()}), textSectionDeleted)    //nolint:exhaustruct
	f.rejected(f.task.UpdateTask(f.ctx, taskModel.UpdateTaskReq{Necessary: f.necessary, ID: taskID, PriorityID: &priorityID, UpdatedAt: f.next()}), textPriorityDeleted) //nolint:exhaustruct
}

// TestLinkToDeletedTaskOrTag — нельзя привязать удаленный тег или тег к удаленной задаче
func TestLinkToDeletedTaskOrTag(t *testing.T) {
	f := newFixture(t)
	liveTaskID := f.createTask(f.inboxID, nil, nil)
	deletedTaskID := f.createTask(f.inboxID, nil, nil)
	f.ok(f.deleteTask(deletedTaskID))
	liveTagID := f.createTag()
	deletedTagID := f.createTag()
	f.ok(f.deleteTag(deletedTagID))

	f.rejected(f.linkTag(deletedTaskID, liveTagID), textTaskDeleted)
	f.rejected(f.linkTag(liveTaskID, deletedTagID), textTagDeleted)
}

// TestBlockerOnDeletedTask — нельзя создать блокер на удаленную задачу ни с одной из сторон
func TestBlockerOnDeletedTask(t *testing.T) {
	f := newFixture(t)
	liveTaskID := f.createTask(f.inboxID, nil, nil)
	deletedTaskID := f.createTask(f.inboxID, nil, nil)
	f.ok(f.deleteTask(deletedTaskID))

	f.rejected(f.createBlocker(uuid.New(), liveTaskID, deletedTaskID), textTaskDeleted)
	f.rejected(f.createBlocker(uuid.New(), deletedTaskID, liveTaskID), textTaskDeleted)
}

// TestDeletedToDeletedAllowed — удаленная строка может ссылаться на удаленную: после удаления задачи
// ее удаленная связь с тегом и удаленный блокер не мешают удалить тег, список, раздел и приоритет
func TestDeletedToDeletedAllowed(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	sectionID := f.createSection(listID)
	priorityID := f.createPriority()
	taskID := f.createTask(listID, &sectionID, &priorityID)
	otherTaskID := f.createTask(listID, nil, nil)
	tagID := f.createTag()
	f.ok(f.linkTag(taskID, tagID))
	blockerID := uuid.New()
	f.ok(f.createBlocker(blockerID, taskID, otherTaskID))

	// Удаляем по правилу: сначала дети, потом родители
	f.ok(f.unlinkTag(taskID, tagID))
	f.ok(f.deleteBlocker(blockerID))
	f.ok(f.deleteTask(taskID))
	f.ok(f.deleteTask(otherTaskID))
	f.ok(f.deleteTag(tagID))
	f.ok(f.deleteSection(sectionID))
	f.ok(f.deletePriority(priorityID))
	f.ok(f.deleteList(listID))

	// Удаленная задача по-прежнему ссылается на удаленные список, раздел и приоритет
	n := f.count(sq.Select("COUNT(*)").From("tasks.tasks").Where(sq.Eq{"id": taskID, "list_id": listID, "section_id": sectionID, "priority_id": priorityID}).Where(sq.NotEq{"deleted_at": nil}))
	if n != 1 {
		t.Fatalf("ссылки удаленной задачи изменились")
	}
}

// TestLastWriteWins — правка старее сохраненной игнорируется без ошибки
func TestLastWriteWins(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	older := f.next()
	newer := f.next()

	newName := "Новое"
	oldName := "Старое"
	f.ok(f.taskList.UpdateTaskList(f.ctx, taskListModel.UpdateTaskListReq{Necessary: f.necessary, ID: listID, Name: &newName, UpdatedAt: newer})) //nolint:exhaustruct
	f.ok(f.taskList.UpdateTaskList(f.ctx, taskListModel.UpdateTaskListReq{Necessary: f.necessary, ID: listID, Name: &oldName, UpdatedAt: older})) //nolint:exhaustruct

	if n := f.count(sq.Select("COUNT(*)").From("tasks.task_lists").Where(sq.Eq{"id": listID, "name": newName})); n != 1 {
		t.Fatalf("старая правка перезаписала новую")
	}

	// Удаление старее последней правки тоже игнорируется
	f.ok(f.taskList.DeleteTaskList(f.ctx, taskListModel.DeleteTaskListReq{Necessary: f.necessary, ID: listID, UpdatedAt: older}))
	if n := f.count(sq.Select("COUNT(*)").From("tasks.task_lists").Where(sq.Eq{"id": listID, "deleted_at": nil})); n != 1 {
		t.Fatalf("старое удаление применилось")
	}
}

// TestIdempotentCreate — повтор создания из офлайн-очереди не создает дубль и не пишет в аудит-лог
func TestIdempotentCreate(t *testing.T) {
	f := newFixture(t)
	taskID := uuid.New()
	request := f.taskReq(taskID, f.inboxID, nil, nil)
	f.ok(f.task.CreateTask(f.ctx, request))
	f.ok(f.task.CreateTask(f.ctx, request))

	if n := f.count(sq.Select("COUNT(*)").From("tasks.tasks").Where(sq.Eq{"id": taskID})); n != 1 {
		t.Fatalf("задач с идентификатором: %d", n)
	}
	if n := f.count(sq.Select("COUNT(*)").From("tasks.audit_log").Where(sq.Eq{"entity_id": taskID.String()})); n != 1 {
		t.Fatalf("записей аудит-лога о создании: %d", n)
	}
}

// TestAuditLog — каждая принятая мутация пишет запись аудит-лога, отклоненная — не пишет
func TestAuditLog(t *testing.T) {
	f := newFixture(t)
	listID := f.createList()
	taskID := f.createTask(listID, nil, nil)
	f.rejected(f.deleteList(listID), textListHasTasks)
	f.ok(f.deleteTask(taskID))
	f.ok(f.deleteList(listID))

	byMethod := func(entityID uuid.UUID, method string) int {
		return f.count(sq.Select("COUNT(*)").From("tasks.audit_log").Where(sq.Eq{"entity_id": entityID.String(), "method": method, "user_id": f.necessary.UserID}))
	}
	if byMethod(listID, "create") != 1 || byMethod(listID, "delete") != 1 {
		t.Fatalf("неверные записи аудит-лога списка")
	}
	if byMethod(taskID, "create") != 1 || byMethod(taskID, "delete") != 1 {
		t.Fatalf("неверные записи аудит-лога задачи")
	}

	// Слепок «до» удаления сохранен
	if n := f.count(sq.Select("COUNT(*)").From("tasks.audit_log").Where(sq.Eq{"entity_id": listID.String(), "method": "delete"}).Where("snapshot_before IS NOT NULL AND snapshot_after IS NULL")); n != 1 {
		t.Fatalf("у удаления нет слепка «до»")
	}
}
