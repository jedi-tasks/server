package integrationTest

import (
	"github.com/google/uuid"

	"server/internal/enum/sortMode"
	blockerModel "server/internal/modules/blocker/model"
	listSectionModel "server/internal/modules/listSection/model"
	priorityModel "server/internal/modules/priority/model"
	tagModel "server/internal/modules/tag/model"
	taskModel "server/internal/modules/task/model"
	taskListModel "server/internal/modules/taskList/model"
)

// createList создает список и возвращает его идентификатор
func (f *fixture) createList() uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.ok(f.taskList.CreateTaskList(f.ctx, taskListModel.CreateTaskListReq{
		Necessary: f.necessary, ID: id, Name: "Список", ColorHex: "#000000", Icon: "list.bullet",
		SortOrder: 1, SortMode: sortMode.Manual, CreatedAt: f.next(),
	}))
	return id
}

// createSection создает раздел в списке
func (f *fixture) createSection(listID uuid.UUID) uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.ok(f.listSection.CreateListSection(f.ctx, f.sectionReq(id, listID)))
	return id
}

// sectionReq собирает запрос на создание раздела
func (f *fixture) sectionReq(id, listID uuid.UUID) listSectionModel.CreateListSectionReq {
	return listSectionModel.CreateListSectionReq{Necessary: f.necessary, ID: id, ListID: listID, Name: "Раздел", SortOrder: 1, UpdatedAt: f.next()}
}

// createPriority создает приоритет
func (f *fixture) createPriority() uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.ok(f.priority.CreatePriority(f.ctx, priorityModel.CreatePriorityReq{
		Necessary: f.necessary, ID: id, Name: "Важно", ColorHex: "#FF0000", Marks: "!!!", SortOrder: 0, UpdatedAt: f.next(),
	}))
	return id
}

// taskReq собирает запрос на создание задачи
func (f *fixture) taskReq(id, listID uuid.UUID, sectionID, priorityID *uuid.UUID) taskModel.CreateTaskReq {
	return taskModel.CreateTaskReq{ //nolint:exhaustruct
		Necessary: f.necessary, ID: id, ListID: listID, SectionID: sectionID, PriorityID: priorityID,
		Title: "Задача", CreatedAt: f.next(),
	}
}

// createTask создает задачу
func (f *fixture) createTask(listID uuid.UUID, sectionID, priorityID *uuid.UUID) uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.ok(f.task.CreateTask(f.ctx, f.taskReq(id, listID, sectionID, priorityID)))
	return id
}

// deleteTask удаляет задачу
func (f *fixture) deleteTask(id uuid.UUID) error {
	return f.task.DeleteTask(f.ctx, taskModel.DeleteTaskReq{Necessary: f.necessary, ID: id, UpdatedAt: f.next()})
}

// deleteList удаляет список
func (f *fixture) deleteList(id uuid.UUID) error {
	return f.taskList.DeleteTaskList(f.ctx, taskListModel.DeleteTaskListReq{Necessary: f.necessary, ID: id, UpdatedAt: f.next()})
}

// deleteSection удаляет раздел
func (f *fixture) deleteSection(id uuid.UUID) error {
	return f.listSection.DeleteListSection(f.ctx, listSectionModel.DeleteListSectionReq{Necessary: f.necessary, ID: id, UpdatedAt: f.next()})
}

// deletePriority удаляет приоритет
func (f *fixture) deletePriority(id uuid.UUID) error {
	return f.priority.DeletePriority(f.ctx, priorityModel.DeletePriorityReq{Necessary: f.necessary, ID: id, UpdatedAt: f.next()})
}

// createTag создает тег
func (f *fixture) createTag() uuid.UUID {
	f.t.Helper()
	id := uuid.New()
	f.ok(f.tag.CreateTag(f.ctx, tagModel.CreateTagReq{Necessary: f.necessary, ID: id, Name: "тег", UpdatedAt: f.next()}))
	return id
}

// deleteTag удаляет тег
func (f *fixture) deleteTag(id uuid.UUID) error {
	return f.tag.DeleteTag(f.ctx, tagModel.DeleteTagReq{Necessary: f.necessary, ID: id, UpdatedAt: f.next()})
}

// linkTag привязывает тег к задаче
func (f *fixture) linkTag(taskID, tagID uuid.UUID) error {
	return f.tag.AddTagToTask(f.ctx, tagModel.LinkTagReq{Necessary: f.necessary, TaskID: taskID, TagID: tagID, UpdatedAt: f.next()})
}

// unlinkTag отвязывает тег от задачи
func (f *fixture) unlinkTag(taskID, tagID uuid.UUID) error {
	return f.tag.RemoveTagFromTask(f.ctx, tagModel.LinkTagReq{Necessary: f.necessary, TaskID: taskID, TagID: tagID, UpdatedAt: f.next()})
}

// createBlocker создает блокер «blockingTaskID блокирует taskID»
func (f *fixture) createBlocker(id, taskID, blockingTaskID uuid.UUID) error {
	return f.blocker.CreateBlocker(f.ctx, blockerModel.CreateBlockerReq{
		Necessary: f.necessary, ID: id, TaskID: taskID, BlockingTaskID: &blockingTaskID, Reason: nil,
		SinceDay: "2026-10-29", UpdatedAt: f.next(),
	})
}

// deleteBlocker удаляет блокер
func (f *fixture) deleteBlocker(id uuid.UUID) error {
	return f.blocker.DeleteBlocker(f.ctx, blockerModel.DeleteBlockerReq{Necessary: f.necessary, ID: id, UpdatedAt: f.next()})
}
