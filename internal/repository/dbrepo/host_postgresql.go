package dbrepo

import (
	"context"
	"github.com/djedjethai/vigilate/internal/models"
	"log"
	"time"
)

// InsertHost insert a host into database
func (m *postgresDBRepo) InsertHost(h models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into hosts(
			host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at
		) 
		values(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		returning id`

	var newID int
	err := m.DB.QueryRowContext(ctx, query,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		h.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		log.Println(err)
		return newID, err
	}
	return newID, nil
}

func (m *postgresDBRepo) GetHostByID(id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var h models.Host
	query := `
		select id, host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at 
		from hosts 
		where id = $1
		`

	rows := m.DB.QueryRowContext(ctx, query, id)

	err := rows.Scan(
		&h.ID,
		&h.HostName,
		&h.CanonicalName,
		&h.URL,
		&h.IP,
		&h.IPV6,
		&h.Location,
		&h.OS,
		&h.Active,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		return h, err
	}

	return h, nil
}

func (m *postgresDBRepo) UpdateHost(h models.Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update hosts set host_name = $1, canonical_name = $2, url = $3, ip = $4, ipv6 = $5, os = $6, active = $7, location = $8, updated_at = $9 where id = $10`

	_, err := m.DB.ExecContext(ctx, query,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.OS,
		h.Active,
		h.Location,
		time.Now(),
		h.ID,
	)

	if err != nil {
		log.Println(err)
		return error(err)
	}
	return nil

}
