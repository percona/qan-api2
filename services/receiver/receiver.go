// qan-api
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

package receiver

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	qanpb "github.com/percona/pmm/api/qan"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Percona-Lab/qan-api/models"
)

// Service implements gRPC service to communicate with agent.
type Service struct {
	mbm models.MetricsBucket
	l   *logrus.Entry
}

// NewService create new insstance of Service.
func NewService(mbm models.MetricsBucket) *Service {
	return &Service{
		mbm: mbm,
		l:   logrus.WithField("component", "receiver"),
	}
}

// DataInterchange implements rpc to exchange data between API and agent.
func (s *Service) DataInterchange(stream qanpb.Agent_DataInterchangeServer) error {
	fmt.Println("Start...")
	for {
		agentMsg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("recved from agent: %+v", agentMsg)
		}
		err = s.mbm.Save(agentMsg)
		if err != nil {
			fmt.Printf("save error: %v \n", err)
			return fmt.Errorf("save error: %v", err)
		}
		savedAmount := len(agentMsg.MetricsBucket)
		fmt.Printf("Rcvd and saved %v Metrics Buckets\n", savedAmount)
		// look for msgs to be sent to client
		msg := qanpb.ApiMessage{SavedAmount: uint32(savedAmount)}
		if err := stream.Send(&msg); err != nil {
			return err
		}
	}
}

func (s *Service) TODO(ctx context.Context, req *qanpb.AgentMessageTODO) (*qanpb.ApiMessageTODO, error) {
	switch t := req.Data.GetTypeUrl(); t {
	case qanpb.AgentMessageTypeURL:
		var msg qanpb.AgentMessage
		if err := proto.Unmarshal(req.Data.Value, &msg); err != nil {
			s.l.Error(err)
			return nil, err
		}

		if err := s.mbm.Save(&msg); err != nil {
			s.l.Error(err)
			return nil, err
		}
		return new(qanpb.ApiMessageTODO), nil

	default:
		s.l.Errorf("Unexpected message type %q.", t)
		return nil, status.Errorf(codes.InvalidArgument, "unexpected message type %q", t)
	}
}
