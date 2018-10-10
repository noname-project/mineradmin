package dbstore

import (
	"github.com/boomstarternetwork/mineradmin/dbstore/model"
	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/getlantern/errors"
)

func (s DBStore) ProjectGet(projectID string) (store.Project, bool, error) {
	project := &model.Project{}

	found, err := s.xdb.ID(projectID).Get(project)
	if err != nil || !found {
		return store.Project{}, found, err
	}

	return store.Project{
		ID:   project.ID,
		Name: project.Name,
	}, true, err
}

func (s DBStore) ProjectAdd(projectID string, projectName string) error {
	_, err := s.xdb.Insert(model.Project{
		ID:   projectID,
		Name: projectName,
	})
	return err
}

func (s DBStore) ProjectSet(projectID string, newProjectName string) error {
	updated, err := s.xdb.ID(projectID).Update(model.Project{
		Name: newProjectName,
	})
	if err != nil {
		return err
	}
	if updated == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (s DBStore) ProjectRemove(projectID string) error {
	_, err := s.xdb.ID(projectID).Delete(model.Project{})
	return err
}
