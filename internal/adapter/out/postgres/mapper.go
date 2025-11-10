package postgres

import (
	"database/sql"
	"time"
)

func nullTimeFromSQL(st sql.NullTime) *time.Time {
	if !st.Valid {
		return nil
	}
	return &st.Time
}

func nullIntFromSQL(si sql.NullInt64) *int64 {
	if !si.Valid {
		return nil
	}
	return &si.Int64
}

func nullIntFromSQL32(si sql.NullInt32) *int32 {
	if !si.Valid {
		return nil
	}
	return &si.Int32
}

func nullFloatFromSQL(sf sql.NullFloat64) *float64 {
	if !sf.Valid {
		return nil
	}
	return &sf.Float64
}

func nullStringFromSQL(ss sql.NullString) *string {
	if !ss.Valid {
		return nil
	}
	return &ss.String
}
