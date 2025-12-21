package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

const iLike = "%%%v%%"

type SortByTypes string

const (
	SortByPublishedAt SortByTypes = "created_at"
	SortByTypesTitle  SortByTypes = "title"
)

type ListNewsParams struct {
	Limit     uint32
	Offset    uint32
	Including string           // ""
	SortBy    SortByTypes      // TITLE
	Ascending bool             // false
	StartTime pgtype.Timestamp // nil
	EndTime   pgtype.Timestamp // nil
}

func (q *Queries) ListNews(ctx context.Context, req ListNewsParams) ([]News, uint32, error) {
	descOrAsc := func() string {
		if req.Ascending {
			return "ASC"
		}
		return "DESC"
	}

	args := []any{}
	queryBuilder := strings.Builder{}

	// beginning of query
	listNews := `
	SELECT id, media_id, creator_username, title, content, created_at, updated_at, deleted_at
	FROM news
	WHERE deleted_at IS NULL
	`
	queryBuilder.WriteString(listNews)

	if req.Including != "" {
		args = append(args, fmt.Sprintf(iLike, req.Including))
		queryBuilder.WriteString(fmt.Sprintf(" AND title ILIKE $%v", len(args)))

		args = append(args, fmt.Sprintf(iLike, req.Including))
		queryBuilder.WriteString(fmt.Sprintf(" OR content ILIKE $%v ", len(args)))
	}

	if req.StartTime.Valid {
		args = append(args, req.StartTime)
		queryBuilder.WriteString(fmt.Sprintf(" AND created_at > $%v ", len(args)))
	}

	if req.EndTime.Valid {
		args = append(args, req.EndTime)
		queryBuilder.WriteString(fmt.Sprintf(" AND created_at < $%v ", len(args)))
	}

	// intercept here to get total pages
	totalPagesArg := append(args, req.Limit)
	totalPagesQuery := strings.Builder{}
	totalPagesQuery.WriteString(fmt.Sprintf(`
		SELECT CEIL(COUNT(*) * 1.0 / $%v) FROM (%s);
	`, len(totalPagesArg), queryBuilder.String()))

	// "ORDER BY" must be the last
	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s", req.SortBy))
	queryBuilder.WriteString(fmt.Sprintf(" %s ", descOrAsc()))

	// end with pagination
	args = append(args, req.Limit)
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%v ", len(args)))
	args = append(args, req.Offset)
	queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%v;", len(args)))

	rows, err := q.db.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []News{}
	for rows.Next() {
		var i News
		if err := rows.Scan(
			&i.ID,
			&i.MediaID,
			&i.CreatorUsername,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	row := q.db.QueryRow(ctx, totalPagesQuery.String(), totalPagesArg...)
	var tot uint32
	err = row.Scan(
		&tot,
	)
	if err != nil {
		return nil, 0, err
	}

	return items, tot, nil
}
