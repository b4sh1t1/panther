package awsglue

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/pantherlog/rewrite_fields"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/pantherlog/tcodec"
)

// TODO: [awsglue] Add more mappings of invalid Athena field name characters here
// NOTE: The mapping should be easy to remember (so no ASCII code etc) and complex enough
// to avoid possible conflicts with other fields.
var fieldNameReplacer = strings.NewReplacer(
	"@", "_at_sign_",
	",", "_comma_",
	"`", "_backtick_",
	"'", "_apostrophe_",
)

func RewriteFieldName(name string) string {
	result := fieldNameReplacer.Replace(name)
	if result == name {
		return name
	}
	return strings.Trim(result, "_")
}

const (
	// We want our output JSON timestamps to be: YYYY-MM-DD HH:MM:SS.fffffffff
	// https://aws.amazon.com/premiumsupport/knowledge-center/query-table-athena-timestamp-empty/
	TimestampLayout     = `2006-01-02 15:04:05.000000000`
	TimestampLayoutJSON = `"` + TimestampLayout + `"`
)

func NewTimestampEncoder() tcodec.TimeEncoder {
	return &timestampEncoder{}
}

var _ tcodec.TimeEncoder = (*timestampEncoder)(nil)

type timestampEncoder struct{}

func (*timestampEncoder) EncodeTime(tm time.Time, stream *jsoniter.Stream) {
	if tm.IsZero() {
		stream.WriteNil()
		return
	}
	buf := stream.Buffer()
	buf = tm.UTC().AppendFormat(buf, TimestampLayoutJSON)
	stream.SetBuffer(buf)
}

func RegisterExtensions(api jsoniter.API) jsoniter.API {
	api.RegisterExtension(rewrite_fields.New(RewriteFieldName))
	api.RegisterExtension(tcodec.NewExtension(tcodec.Config{
		// Force all timestamps to be awsglue format and UTC. This is needed to be able to write
		DefaultCodec: tcodec.Join(tcodec.StdCodec(), NewTimestampEncoder()),
		DecorateCodec: func(codec tcodec.TimeCodec) tcodec.TimeCodec {
			dec, _ := tcodec.Split(codec)
			enc := NewTimestampEncoder()
			return tcodec.Join(dec, enc)
		},
	}))
	return api
}