package stream

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"

	pb "github.com/AtlasInsideCorp/UTMStackAgent/agent"
	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"github.com/quantfall/holmes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var IRMutex sync.Mutex

// IncidentResponse runs in the background and periodically checks for commands to execute from the server.
// If commands are found, it executes them and sends the response back to the server.
func StartStream(cnf configuration.Config, client pb.AgentServiceClient, ctx context.Context, cancel context.CancelFunc, h *holmes.Logger) {
	path, err := utils.GetMyPath()
	if err != nil {
		fmt.Printf("Failed to get current path: %v", err)
		h.FatalError("failed to get current path: %v", err)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "token", cnf.AgentToken)
	ctx = metadata.AppendToOutgoingContext(ctx, "agent-id", strconv.Itoa(int(cnf.AgentID)))

	stream, err := client.AgentStream(ctx)
	if err != nil {
		fmt.Printf("Failed to start AgentStream: %v", err)
		h.FatalError("failed to start AgentStream: %v", err)
	}

	// Send the authentication response
	authResponse := &pb.BidirectionalStream{
		StreamMessage: &pb.BidirectionalStream_AuthResponse{
			AuthResponse: &pb.AuthResponse{
				Token:   cnf.AgentToken,
				AgentId: uint64(cnf.AgentID),
			},
		},
	}
	if err := stream.Send(authResponse); err != nil {
		fmt.Printf("Failed to send AuthResponse: %v", err)
		h.FatalError("failed to send AuthResponse: %v", err)
	}

	// Handle the bidirectional stream
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				cancel()
				return
			}
			if err != nil {
				fmt.Printf("Failed to receive message from server: %v", err)
				h.FatalError("failed to receive message from server: %v", err)
			}
			switch msg := in.StreamMessage.(type) {
			case *pb.BidirectionalStream_Command:
				var result string
				// Handle the received command (replace with your logic)
				fmt.Printf("Received command: %s", msg.Command.Command)
				h.Info("received command: %s", msg.Command.Command)
				switch runtime.GOOS {
				case "windows":
					result, _ = utils.Execute("cmd.exe", path, "/C", msg.Command.Command)
				case "linux":
					result, _ = utils.Execute("sh", path, "-c", msg.Command.Command)
				default:
					result, _ = fmt.Sprintf("unsupported operating system: %s", runtime.GOOS), true
				}
				// Send the result back to the server
				if err := stream.Send(&pb.BidirectionalStream{
					StreamMessage: &pb.BidirectionalStream_Result{
						Result: &pb.CommandResult{Result: result, Token: cnf.AgentToken, ExecutedAt: timestamppb.Now(), CmdId: in.GetCommand().CmdId},
					},
				}); err != nil {
					fmt.Printf("Failed to send result to server: %v", err)
					h.FatalError("failed to send result to server: %v", err)
				}
			}
		}
	}()
}
