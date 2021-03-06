package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-yaml/yaml"
	"github.com/gotti/cevigo/pkg/cevioai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/gotti/cevigo/spec"
)

type ttsServer struct {
	talker cevioai.ITalker2V40
	pb.UnimplementedTtsServer
	tokens tokens
	mtx    sync.Mutex
}

func (s *ttsServer) applyParameters(p *pb.CevioTtsRequest) error {
	s.talker.SetCast(p.Cast)
	s.talker.SetVolume(int(p.Volume))
	s.talker.SetSpeed(int(p.Speed))
	s.talker.SetTone(int(p.Pitch))      //高さ
	s.talker.SetAlpha(int(p.Alpha))     //声質
	s.talker.SetToneScale(int(p.Intro)) //抑揚
	comp, err := s.talker.GetComponents()
	if err != nil {
		return fmt.Errorf("getting components: %w", err)
	}
	for k, p := range p.Emotions {
		n, err := comp.ByName(k)
		if err != nil {
			return fmt.Errorf("getting single component: %w", err)
		}
		if err := n.SetValue(int(p)); err != nil {
			return fmt.Errorf("setting value to component: %w", err)
		}
	}
	return nil
}

func (s *ttsServer) speak(text string) ([]byte, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("generating rand: %w", err)
	}
	fPath := fmt.Sprintf("%8x", buf)
	fPath = filepath.Join(filepath.Join(os.TempDir(), "cevigo"), fPath)
	err = os.MkdirAll(filepath.Dir(fPath), os.FileMode(0755))
	if err != nil {
		return nil, fmt.Errorf("making dir: %w", err)
	}
	b, err := s.talker.OutputWaveToFile(text, fPath)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, fmt.Errorf("failed to outputting wav, please check the packet you sent")
	}
	defer os.Remove(fPath)
	if err != nil {
		return nil, fmt.Errorf("outputting: %w", err)
	}
	f, err := os.Open(fPath)
	audio, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return audio, nil
}

func (s *ttsServer) CreateWav(ctx context.Context, req *pb.CevioTtsRequest) (*pb.CevioTtsResponse, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if err := s.validateArgument(req); err != nil {
		log.Printf("validating: %v", err)
		return nil, err
	}
	if err := s.applyParameters(req); err != nil {
		log.Printf("applying parameters: %v", err)
		return nil, err
	}
	audio, err := s.speak(req.Text)
	if err != nil {
		log.Printf("speaking: %v", err)
		return nil, err
	}
	log.Printf("successfully outputted: %s\n", req.Text)
	return &pb.CevioTtsResponse{Audio: audio}, nil
}

func (s *ttsServer) validateArgument(req *pb.CevioTtsRequest) error {
	if !validateToken(s.tokens, req.Token) {
		return status.Errorf(codes.PermissionDenied, "invalid token")
	}
	casts, err := s.talker.GetAvailableCasts()
	if err != nil {
		return fmt.Errorf("get available casts: %w", err)
	}
	arr, err := casts.ToGoArray()
	if !inArray(arr, req.Cast) {
		return status.Errorf(codes.InvalidArgument, "invalid cast")
	}
	if err != nil {
		return fmt.Errorf("available casts to go array: %w", err)
	}
	if req.Volume > 100 {
		return status.Errorf(codes.InvalidArgument, "invalid volume")
	}
	if req.Speed > 100 {
		return status.Errorf(codes.InvalidArgument, "invalid speed")
	}
	if req.Pitch > 100 {
		return status.Errorf(codes.InvalidArgument, "invalid pitch")
	}
	if req.Alpha > 100 {
		return status.Errorf(codes.InvalidArgument, "invalid alpha")
	}
	if req.Intro > 100 {
		return status.Errorf(codes.InvalidArgument, "invalid intro")
	}
	//TODO: validate emotions before processing
	return nil
}

func inArray(array []string, data string) bool {
	for _, a := range array {
		if a == data {
			return true
		}
	}
	return false
}

func validateToken(ts tokens, t string) bool {
	for _, ti := range ts {
		if ti == t {
			return true
		}
	}
	return false
}

func main() {
	apidiff := flag.String("api", "cevio", "cevio, or cevioai")
	flag.Parse()
	var apiname string
	if *apidiff == "cevio" {
		apiname = cevioai.CevioApiName
	} else if *apidiff == "cevioai" {
		apiname = cevioai.CevioAiApiName
	} else {
		log.Println("set cevio or cevioai to --api")
		os.Exit(1)
	}

	t, err := loadTokens("./tokens.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	talker := cevioai.NewITalker2V40(apiname)
	log.Printf("connected to %s\n", apiname)

	lis, err := net.Listen("tcp", ":10001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTtsServer(s, &ttsServer{talker: talker, tokens: t})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if err != nil {
		log.Fatal(err)
	}
}

type tokens []string

func loadTokens(path string) (tokens, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ts tokens
	if err := yaml.Unmarshal(buf, &ts); err != nil {
		return nil, err
	}
	return ts, nil
}
