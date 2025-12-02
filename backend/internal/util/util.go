package util

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// protoのtimestamppb.Timestampをtime.Timeに変換
func ParseTimestampToTime(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}
	t := timestamp.AsTime()
	return &t
}

// time.Timeをpsgoreのtimestamptzに変換
func PaeseTimeToPgtypeTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{
			Valid: false,
		}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

// GoのtimeをPgのtimestamptzに変換
func GoTimeToPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{
			Valid: false,
		}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

// PgのtimestamptzをGoのtimeに変換
func PgTimestamptzToGoTime(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
}
