package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	Username    string `json:"username"`
	Seconds     int    `json:"seconds"`
	TopLanguage string `json:"topLanguage"`
	IsYou       bool   `json:"isYou"`
}

type LeaderboardResult struct {
	Entries  []LeaderboardEntry `json:"entries"`
	YourRank *int               `json:"yourRank"`
}

// GetLeaderboard returns ranked users by total coding time
// within the given range, respecting privacy settings.
func GetLeaderboard(ctx context.Context, pool *pgxpool.Pool, rangeSQL string, limit int, currentUserID string) (*LeaderboardResult, error) {
	rows, err := pool.Query(ctx, `
		SELECT
			u.id,
			u.username,
			SUM(s.duration_seconds) as total_seconds,
			(
				SELECT language FROM sessions s2
				WHERE s2.user_id = u.id
				GROUP BY language
				ORDER BY SUM(duration_seconds) DESC
				LIMIT 1
			) as top_language
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN privacy_settings p ON p.user_id = u.id
		WHERE `+rangeSQL+`
			AND (p.hide_leaderboard IS NULL OR p.hide_leaderboard = false)
			AND (p.profile_public IS NULL OR p.profile_public = true)
			AND u.deleted_at IS NULL
		GROUP BY u.id, u.username
		ORDER BY total_seconds DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allEntries []LeaderboardEntry
	rank := 1
	var yourRank *int

	for rows.Next() {
		var e LeaderboardEntry
		var userID string
		var topLang *string
		if err := rows.Scan(&userID, &e.Username, &e.Seconds, &topLang); err != nil {
			return nil, err
		}
		if topLang != nil {
			e.TopLanguage = *topLang
		}
		e.Rank = rank
		e.IsYou = userID == currentUserID

		if e.IsYou {
			r := rank
			yourRank = &r
		}

		allEntries = append(allEntries, e)
		rank++
	}

	result := &LeaderboardResult{YourRank: yourRank}
	if len(allEntries) > limit {
		result.Entries = allEntries[:limit]
	} else {
		result.Entries = allEntries
	}

	return result, nil
}
