package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func TextToString(t pgtype.Text) string {
	return t.String
}

func TextToStringPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func NumericToFloat32(n pgtype.Numeric) float32 {
	if !n.Valid {
		return 0
	}
	f, _ := n.Float64Value()
	return float32(f.Float64)
}

func Int4ToInt32(i pgtype.Int4) int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

func Float32ToNumeric(f float32) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%f", f))
	return n
}
