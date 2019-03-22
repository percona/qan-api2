// qan-api2
// Copyright (C) 2019 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package analitycs

import (
	"context"
	"errors"
	"time"

	"github.com/percona/pmm/api/qanpb"
)

// GetFilters implements rpc to get list of available labels.
func (s *Service) GetFilters(ctx context.Context, in *qanpb.FiltersRequest) (*qanpb.FiltersReply, error) {
	from := time.Unix(in.GetFrom().Seconds, 0)
	to := time.Unix(in.GetTo().Seconds, 0)
	if from.After(to) {
		return &qanpb.FiltersReply{}, errors.New("from-date cannot be bigger then to-date")
	}
	return s.rm.SelectFilters(from, to)
}
