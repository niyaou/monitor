package main1

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	pb "desay.com/radar-monitor/pb"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 50051, "The server port")
)

type monitorViewServer struct {
	pb.UnimplementedMonitorViewServer
	mu sync.Mutex // protects routeNotes
}

// GetFeature returns the feature at the given point.
func (s *monitorViewServer) Fc2ParamCfg(ctx context.Context, point *pb.RadarParamCfgType) (*pb.RadarParamCfgType, error) {
	// for _, feature := range s.savedFeatures {
	// 	if proto.Equal(feature.Location, point) {
	// 		return feature, nil
	// 	}
	// }
	var msg *pb.RadarParamCfgType = &pb.RadarParamCfgType{}
	msg.U16Cmd = 1
	msg.U16AcqNrFrames = 1

	var payload *pb.TEF82Xx_FrameOptionalParamType = &pb.TEF82Xx_FrameOptionalParamType{}
	payload.U32SeqInterval = 1000
	msg.StrTef82XxFrameOpParam = payload

	// } else {
	// 	var payload *pb.RFE_ChirpShapeType = &pb.RFE_ChirpShapeType{}
	// 	payload.U32TStart = 1000
	// 	payload.U32TPreSampling = 500
	// 	payload.U32TPostSampling = 100
	// 	payload.U32TReturn = 10000
	// 	payload.U32CenterFrequency = 76500
	// 	payload.U32AcqBandwidth = 1000
	// 	payload.U8TxChannelEnable = 63
	// 	payload.U32ATxChannelPower = []uint32{1001, 1002, 1003, 1004, 1005, 1006}
	// 	payload.U8ARxChannelGain = []uint32{71, 72, 73, 74, 75, 76, 77, 78}

	// 	data, _ := proto.Marshal(payload)

	// 	msg.PMsgCbk = data
	// 	src, _ := proto.Marshal(msg)
	// 	encodedStr := hex.EncodeToString(src)
	// 	data1, _ := hex.DecodeString(encodedStr)
	// 	for _, v := range data1 {
	// 		fmt.Printf("%#X \n", v)
	// 	}
	// 	// fmt.Println(src)
	// 	// 48656c6c6f -> 48(4*16+8=72) 65(6*16+5=101) 6c 6c 6f
	// 	// fmt.Println(encodedStr)
	// }

	return msg, nil
}

func newServer() *monitorViewServer {
	s := &monitorViewServer{}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMonitorViewServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
