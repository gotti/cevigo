package cevioai

import (
	"fmt"

	"github.com/go-ole/go-ole"
)

type ISpeakingState2 struct {
	dispatch *ole.IDispatch
}

func NewISpeakingState2(dispatch *ole.IDispatch) ISpeakingState2 {
	return ISpeakingState2{dispatch: dispatch}
}

func (s *ISpeakingState2) IsCompleted() (bool, error) {
	v, err := s.dispatch.GetProperty("IsCompleted")
	if err != nil {
		return false, err
	}
	i, ok := v.Value().(bool)
	if !ok {
		return false, fmt.Errorf("Failed to interface assertion in iscompleted")
	}
	return i, nil
}

func (s *ISpeakingState2) IsSucceeded() (bool, error) {
	v, err := s.dispatch.GetProperty("IsSucceeded")
	if err != nil {
		return false, err
	}
	i, ok := v.Value().(bool)
	if !ok {
		return false, fmt.Errorf("Failed to interface assertion in issucceeded")
	}
	return i, nil
}

func (s *ISpeakingState2) Wait() error {
	_, err := s.dispatch.CallMethod("Wait")
	if err != nil {
		return err
	}
	return nil
}

func (s *ISpeakingState2) Wait_2() error {
	_, err := s.dispatch.CallMethod("Wait_2")
	if err != nil {
		return err
	}
	return nil
}
