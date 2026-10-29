package integrationTest

import (
	"encoding/json"
	"strings"
	"testing"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"

	blockerEndpointMCP "server/internal/modules/blocker/endpoint/mcp"
	tagEndpointMCP "server/internal/modules/tag/endpoint/mcp"
	taskEndpointMCP "server/internal/modules/task/endpoint/mcp"
	taskListEndpointMCP "server/internal/modules/taskList/endpoint/mcp"
	"server/internal/utils/contextKeys"
	mcpTransport "server/internal/utils/mcp"
)

// mcpSession поднимает MCP-сервер с инструментами поверх сервисов фикстуры и подключает к нему
// клиента через in-memory транспорт. Пользователь и устройство — в контексте сервера, как после
// проверки Bearer-токена
func (f *fixture) mcpSession() *sdk.ClientSession {
	f.t.Helper()
	server := mcpTransport.NewServer("test")
	taskEndpointMCP.Register(server, f.task)
	taskListEndpointMCP.Register(server, f.taskList)
	tagEndpointMCP.Register(server, f.tag)
	blockerEndpointMCP.Register(server, f.blocker)

	ctx := contextKeys.SetUserID(f.ctx, f.necessary.UserID)
	ctx = contextKeys.SetDeviceID(ctx, f.necessary.DeviceID)

	serverTransport, clientTransport := sdk.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		f.t.Fatal(err)
	}
	session, err := sdk.NewClient(&sdk.Implementation{Name: "test", Version: "0"}, nil).Connect(ctx, clientTransport, nil) //nolint:exhaustruct
	if err != nil {
		f.t.Fatal(err)
	}
	f.t.Cleanup(func() { _ = session.Close() })
	return session
}

// call вызывает инструмент и разбирает структурированный ответ; ошибка инструмента — в тексте
func (f *fixture) call(session *sdk.ClientSession, name string, args map[string]any, out any) error {
	f.t.Helper()
	result, err := session.CallTool(f.ctx, &sdk.CallToolParams{Name: name, Arguments: args}) //nolint:exhaustruct
	if err != nil {
		f.t.Fatal(err)
	}
	if result.IsError {
		text, _ := result.Content[0].(*sdk.TextContent)
		return &toolError{text: text.Text}
	}
	if out != nil {
		raw, _ := json.Marshal(result.StructuredContent)
		if err := json.Unmarshal(raw, out); err != nil {
			f.t.Fatal(err)
		}
	}
	return nil
}

// toolError — ошибка, которую инструмент вернул агенту текстом
type toolError struct{ text string }

func (e *toolError) Error() string { return e.text }

func TestMCP_ToolsListed(t *testing.T) {
	f := newFixture(t)
	session := f.mcpSession()

	tools, err := session.ListTools(f.ctx, nil)
	f.ok(err)
	names := make([]string, 0, len(tools.Tools))
	for _, tool := range tools.Tools {
		names = append(names, tool.Name)
	}
	if len(names) != 17 {
		t.Fatalf("ожидалось 17 инструментов, получено %d: %v", len(names), names)
	}
}

func TestMCP_TaskLifecycle(t *testing.T) {
	f := newFixture(t)
	session := f.mcpSession()

	// Задача без списка попадает в Inbox
	var created struct{ ID string }
	f.ok(f.call(session, "create_task", map[string]any{"title": "Через MCP", "plannedDay": "2026-11-01"}, &created))

	var tasks struct {
		Tasks []struct {
			ID, ListID, Title string
			PlannedDay        *string
		}
	}
	f.ok(f.call(session, "get_tasks", map[string]any{"ids": []string{created.ID}}, &tasks))
	if len(tasks.Tasks) != 1 || tasks.Tasks[0].ListID != f.inboxID.String() || *tasks.Tasks[0].PlannedDay != "2026-11-01" {
		t.Fatalf("неожиданные задачи: %+v", tasks)
	}

	// Сброс срока пустой строкой. Декодируем в новую переменную: json.Unmarshal переиспользует
	// элементы уже заполненного среза и оставил бы прежний plannedDay
	f.ok(f.call(session, "update_task", map[string]any{"id": created.ID, "plannedDay": ""}, nil))
	var afterReset struct {
		Tasks []struct{ PlannedDay *string }
	}
	f.ok(f.call(session, "get_tasks", map[string]any{"ids": []string{created.ID}}, &afterReset))
	if afterReset.Tasks[0].PlannedDay != nil {
		t.Fatalf("срок не сброшен: %q", *afterReset.Tasks[0].PlannedDay)
	}

	// Удаление с живой связью отклоняется текстом, после снятия тега проходит
	var tag struct{ ID string }
	f.ok(f.call(session, "create_tag", map[string]any{"name": "mcp"}, &tag))
	f.ok(f.call(session, "add_tag_to_task", map[string]any{"taskID": created.ID, "tagID": tag.ID}, nil))
	err := f.call(session, "delete_task", map[string]any{"id": created.ID}, nil)
	if err == nil || !strings.Contains(err.Error(), "live tag links") {
		t.Fatalf("ожидался отказ из-за тега, получено: %v", err)
	}
	f.ok(f.call(session, "remove_tag_from_task", map[string]any{"taskID": created.ID, "tagID": tag.ID}, nil))
	f.ok(f.call(session, "delete_task", map[string]any{"id": created.ID}, nil))

	var afterDelete struct{ Tasks []struct{ ID string } }
	f.ok(f.call(session, "get_tasks", map[string]any{"ids": []string{created.ID}}, &afterDelete))
	if len(afterDelete.Tasks) != 0 {
		t.Fatalf("удаленная задача видна: %+v", afterDelete)
	}
}

func TestMCP_BlockerRequiresOneOfSources(t *testing.T) {
	f := newFixture(t)
	session := f.mcpSession()

	var task struct{ ID string }
	f.ok(f.call(session, "create_task", map[string]any{"title": "Заблокированная"}, &task))

	// Без причины и без блокирующей задачи — отказ
	err := f.call(session, "create_blocker", map[string]any{"taskID": task.ID}, nil)
	if err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("ожидался отказ, получено: %v", err)
	}

	// С причиной — создается и снимается
	var blocker struct{ ID string }
	f.ok(f.call(session, "create_blocker", map[string]any{"taskID": task.ID, "reason": "ждем ответа"}, &blocker))
	f.ok(f.call(session, "resolve_blocker", map[string]any{"id": blocker.ID}, nil))

	var blockers struct {
		Blockers []struct{ ResolvedAt *string }
	}
	f.ok(f.call(session, "get_blockers", map[string]any{"taskIDs": []string{task.ID}}, &blockers))
	if len(blockers.Blockers) != 1 || blockers.Blockers[0].ResolvedAt == nil {
		t.Fatalf("блокер не снят: %+v", blockers)
	}
}
