package inmemory

import (
	"sync"
	"time"

	"github.com/vMaroon/SmartDispenserCounter/datatypes"
	"github.com/vMaroon/SmartDispenserCounter/pkg/db"
	"github.com/vMaroon/SmartDispenserCounter/queries"
)

// newDispenserJobDB returns a new instance of dispenserJobDB.
func newDispenserJobDB() *dispenserJobDB {
	return &dispenserJobDB{
		dispenserJobs: make([]*datatypes.DispenserJob, db.DispenserJobsHistoryLimit),
		currentJobID:  0,
		lock:          sync.Mutex{},
	}
}

// dispenserJobDB struct implements DB logic required for storing and handling dispenser-jobs.
type dispenserJobDB struct {
	dispenserJobs []*datatypes.DispenserJob
	currentJobID  uint64
	lock          sync.Mutex
}

//nolint
// GetDispenserJobs function to return jobs with sorting by timestamp.
// This function is used to fetch history of up to historyLimit latest jobs.
func (dispenserJobDB *dispenserJobDB) GetDispenserJobs(sortable *queries.Sortable) []datatypes.DispenserJob {
	if dispenserJobDB.currentJobID < 0 {
		return []datatypes.DispenserJob{}
	}

	result := make([]datatypes.DispenserJob, 0, db.DispenserJobsHistoryLimit)

	if dispenserJobDB.currentJobID >= db.DispenserJobsHistoryLimit {
		// overpassed limit, array is cyclic - starting point is current index + 1
		for i := 0; i < db.DispenserJobsHistoryLimit; i++ {
			index := (dispenserJobDB.currentJobID + uint64(i) + 1) % db.DispenserJobsHistoryLimit

			if dispenserJobDB.dispenserJobs[index] == nil {
				continue
			}

			result = append(result, *dispenserJobDB.dispenserJobs[index])

			if index == dispenserJobDB.currentJobID%db.DispenserJobsHistoryLimit {
				break
			}
		}
	} else {
		for i := 0; i < db.DispenserJobsHistoryLimit; i++ {
			if dispenserJobDB.dispenserJobs[i] == nil {
				continue
			}

			result = append(result, *dispenserJobDB.dispenserJobs[i])

			if i == int(dispenserJobDB.currentJobID) {
				break
			}
		}
	}

	if sortable.AscendingOrder {
		return result
	}

	return reverseJobsSlice(result)
}

// GetNextJob function to return jobs with sorting by timestamp.
func (dispenserJobDB *dispenserJobDB) GetNextJob() *datatypes.DispenserJob {
	return dispenserJobDB.dispenserJobs[dispenserJobDB.currentJobID%db.DispenserJobsHistoryLimit]
}

// InsertDispenserJob function to insert a job to the database.
// If the ID field is left empty, a unique ID is assigned.
// Otherwise, the given entry may overwrite an existing one.
//
// Returns the job's ID.
func (dispenserJobDB *dispenserJobDB) InsertDispenserJob(dispenserJob *datatypes.DispenserJob) uint64 {
	dispenserJobDB.currentJobID++

	dispenserJob.ID = dispenserJobDB.currentJobID
	dispenserJob.CreationTimestamp = time.Now()

	dispenserJobDB.dispenserJobs[dispenserJobDB.currentJobID%db.DispenserJobsHistoryLimit] = dispenserJob

	return dispenserJob.ID
}

// DeleteDispenserJob function to delete a job given its unique identifier.
// If a relevant job is found and is deleted, returns true. Otherwise, returns false.
func (dispenserJobDB *dispenserJobDB) DeleteDispenserJob(jobID uint64) bool {
	for i := 0; i < len(dispenserJobDB.dispenserJobs); i++ {
		if dispenserJobDB.dispenserJobs[i] != nil && dispenserJobDB.dispenserJobs[i].ID == jobID {
			dispenserJobDB.dispenserJobs[i] = nil
			return true
		}
	}

	return false
}

func reverseJobsSlice(slice []datatypes.DispenserJob) []datatypes.DispenserJob {
	result := make([]datatypes.DispenserJob, 0, len(slice))

	for i := len(slice) - 1; i >= 0; i-- {
		result = append(result, slice[i])
	}

	return result
}
