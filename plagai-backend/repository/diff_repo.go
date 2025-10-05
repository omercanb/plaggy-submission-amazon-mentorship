package repository

import (
	"errors"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"gorm.io/gorm"
)

var ErrDiffDatabase = errors.New("database error")

type DiffRepository interface {
	// Fetch diffs for a student-assignment using start/finish as indices.
	// start: 0-based inclusive offset
	// finish: exclusive upper bound; if finish == -1, fetch everything after start.
	GetDiffs(studentAssignmentID uint, start, finish int) ([]domain.Diff, error)
}

type diffRepository struct {
	db *gorm.DB
}

func NewDiffRepository(db *gorm.DB) DiffRepository {
	return &diffRepository{db: db}
}

func (r *diffRepository) GetDiffs(studentAssignmentID uint, start, finish int) ([]domain.Diff, error) {
	if start < 0 {
		start = 0
	}

	query := r.db.
		Where("student_assignment_id = ?", studentAssignmentID).
		Order("created_at ASC")

	// Apply limit/offset only if finish is non-negative.
	if finish >= 0 {
		if finish < start {
			return []domain.Diff{}, nil
		}
		limit := finish - start
		query = query.Offset(start).Limit(limit)
	} else {
		// finish == -1 ? fetch everything from start onward.
		query = query.Offset(start)
	}

	var dbDiffs []database.Diff
	if err := query.Find(&dbDiffs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []domain.Diff{}, nil
		}
		return nil, ErrDiffDatabase
	}

	diffs := make([]domain.Diff, len(dbDiffs))
	for i, d := range dbDiffs {
		diffs[i] = toDomainDiff(&d)
	}
	return diffs, nil
}

func toDomainDiff(d *database.Diff) domain.Diff {
	return domain.Diff{
		ID:        d.ID,
		FilePath:  d.FilePath,
		PatchText: d.DiffData,
		Timestamp: d.CreatedAt,
	}
}
