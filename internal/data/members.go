package data

import (
	"context"
	"time"
)

type Member struct {
	ID         int       `json:"id"`
	PlayerID   string    `json:"player_id"`
	MemberType string    `json:"member_type"`
	GroupID    int       `json:"group_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MemberReponse struct {
	ID         int       `json:"id"`
	PlayerID   string    `json:"player_id"`
	GroupID    int       `json:"group_id"`
	MemberType string    `json:"member_type"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GET/members/:group_id
func (m *Member) GetAllMemberOfAGroup(group_id int) ([]*MemberReponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		select 
			m.id,
			m.player_id,
			m.group_id,
			m.member_type,
			u.first_name,
			u.last_name,
			u.email,
			m.created_at,
			m.updated_at
		from
			members m
		inner join users u
			on m.player_id = u.id
		where m.group_id = $1
	`

	rows, err := db.QueryContext(ctx, query, group_id)
	if err != nil {
		return nil, err
	}

	var members []*MemberReponse
	for rows.Next() {
		var member MemberReponse
		err := rows.Scan(
			&member.ID,
			&member.PlayerID,
			&member.GroupID,
			&member.MemberType,
			&member.FirstName,
			&member.LastName,
			&member.Email,
			&member.CreatedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

// POST/members/create
func (m *Member) CreateMember(member Member) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var newId int
	query := `
		insert into members (player_id, group_id, member_type, created_at, updated_at)
		values ($1, $2, $3, $4, $5) returning id
	`
	err := db.QueryRowContext(
		ctx,
		query,
		member.PlayerID,
		member.GroupID,
		member.MemberType,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}
