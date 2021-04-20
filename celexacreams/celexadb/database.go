package celexadb

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4/pgxpool"
	"strings"
	"time"
)

var Dbpool *pgxpool.Pool

// SanityCheck checks whether the database is in a usable state
func SanityCheck() error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	// select tablename from pg_tables where schemaname = 'public';
	builder := psql.Select("tablename").From("pg_tables").Where("schemaname = 'public'")
	q, _, _ := builder.ToSql()
	var tables []string
	rows, err := Dbpool.Query(context.Background(), q)
	if err != nil {
		return err
	}
	// Iterate through the result set
	for rows.Next() {
		var n string
		err = rows.Scan(&n)
		if err != nil {
			return err
		}
		tables = append(tables, n)
	}
	if len(tables) < 1 {
		return fmt.Errorf("This database is empty, we can't proceed like this.")
	}
	return nil
}

// GuildCreate updates the guilds table on GuildCreate events
func GuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) error {
	layout := "2006-01-02T15:04:05.999999999Z07:00"
	t, err := time.Parse(layout, string(g.JoinedAt))
	if err != nil {
		return err
	}
	name := strings.Replace(g.Name, "'", "", -1)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Insert("guilds").Columns("id", "name", "join_date", "owner").
		Values(g.ID, name, t, g.OwnerID).
		Suffix(fmt.Sprintf("on conflict (id) do update set name = '%s', owner = '%v'", name, g.OwnerID))
	q, args, _ := builder.ToSql()
	_, err = Dbpool.Exec(context.Background(), q, args...)
	if err != nil {
		return err
	}
	builder = psql.Insert("message_counts").Columns("guild_id").
		Values(g.ID).Suffix(fmt.Sprintf("on conflict (guild_id) do nothing"))
	q, args, _ = builder.ToSql()
	_, err = Dbpool.Exec(context.Background(), q, args...)
	if err != nil {
		return err
	}
	return nil
}

// IncrementMessageCount increments the counter in the message_count table
func IncrementMessageCount(g string) error {
	q := fmt.Sprintf("update message_counts set count = count + 1 where guild_id = %v", g)
	_, err := Dbpool.Exec(context.Background(), q)
	return err
}