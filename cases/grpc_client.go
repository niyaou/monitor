/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"desay.com/radar-monitor/logger"
	pb "desay.com/radar-monitor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/examples/data"
	"google.golang.org/protobuf/proto"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func Acknowledged(client pb.MonitorViewClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	ack := &pb.RadarParamCfgAckType{}
	ack.EnuParamCfgAck = pb.AckCode_ParamCfgType_ERR_BUSY
	stream, _ := client.Acknowledged(ctx, ack)
	logger.Info(" send >>>>>>>>>>> to service>>>>>>>> %v  %T", ack, ack)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			logger.Error(" err %v  %T", err, err)
			return err
		}
		if err != nil {
			logger.Error(" err %v  %T", err, err)
			return err
		}
		logger.Info(" data %v  %T", data, data)
	}
}

func Fc2ParamCfg(client pb.MonitorViewClient, point *pb.RadarParamCfgType) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fc_payload, err := client.Fc2ParamCfg(ctx, point)
	if err != nil {
		log.Fatalf("client.GetFeature failed: %v", err)
	}
	// PC_RsdkTef82XXAutoDriftParamsT

	var msg *pb.TEF82Xx_FrameOptionalParamType = &pb.TEF82Xx_FrameOptionalParamType{}
	// msg = fc_payload.StrTef82XxFrameOpParam
	data, _ := proto.Marshal(fc_payload.StrTef82XxFrameOpParam)
	logger.Info(" data %v  %T", data, data)
	proto.Unmarshal(data, msg)

	logger.Info(" msg %v  %T", msg, msg)
	log.Println(fc_payload)
	// log.Println(msg)

}

func ConfigFc2Param(client pb.MonitorViewClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	point := &pb.RadarParamCfgType{}
	point.U16Cmd = 555
	point.U16AcqNrFrames = 555

	var msg1 *pb.TEF82Xx_FrameOptionalParamType = &pb.TEF82Xx_FrameOptionalParamType{}
	msg1.U32SeqInterval = 99
	point.StrTef82XxFrameOpParam = msg1

	count := 0
	updateTicker := time.NewTicker(1000 * time.Millisecond)
	defer func() {
		updateTicker.Stop()
	}()
	for {
		if count > 60 {
			break
		}
		count++
		<-updateTicker.C
		_, err := client.ConfigFc2Param(ctx, point)
		logger.Info(" ConfigFc2Param------- %v --------- ", 3)
		if err != nil {
			log.Fatalf("client.GetFeature failed: %v", err)
		}
	}

}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewMonitorViewClient(conn)

	// Fc2ParamCfg(client, &pb.RadarParamCfgType{})
	// ConfigFc2Param(client)
	Acknowledged(client)

}
