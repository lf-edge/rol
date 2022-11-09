package tests

import (
	"github.com/google/uuid"
	"testing"
)

var (
	repoTesters = []IGenericRepositoryTester{
		NewGormSQLiteGenericRepositoryTester[uuid.UUID, TestEntityUUID](),
		NewGormSQLiteGenericRepositoryTester[int, TestEntityInt](),
	}
)

func Test_GenericRepo_InsertAndDelete(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestInsertAndDelete()
	}
}

func Test_GenericRepo_DeleteAll(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestDeleteAll()
	}
}

func Test_GenericRepo_GetByID(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestGetByID()
	}
}

func Test_GenericRepo_GetByIDExtended(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestGetByIDExtended()
	}
}

func Test_GenericRepo_Delete(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestDelete()
	}
}

func Test_GenericRepo_IsExist(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestIsExist()
	}
}

func Test_GenericRepo_Count_All(t *testing.T) {
	for _, tester := range repoTesters {
		tester.SetTesting(t)
		tester.TestCount()
	}
}

func Test_GenericRepo_Clean(t *testing.T) {
	for _, tester := range repoTesters {
		err := tester.Dispose()
		if err != nil {
			t.Error("Failed to clean repository tester stuff")
		}
	}
}
