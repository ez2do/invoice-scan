package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"invoice-scan/backend/pkg/log"
	"sync"
	"testing"
	"time"
)

func TestErrGroupPanicRecovery(t *testing.T) {
	g := NewErrGroupWithRecovery(context.Background())

	g.Go(func() error {
		panic("a func has panic")
	})

	if err := g.Wait(); err != nil {
		assert.Contains(t, err.Error(), "a func has panic")
	} else {
		t.Fatal("Expected panic error recovered")
	}
}

func TestErrGroupContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	g := NewErrGroupWithRecovery(ctx)
	g.Go(func() error {
		time.Sleep(3 * time.Second)
		return nil
	})

	err := g.Wait()
	assert.EqualError(t, err, context.DeadlineExceeded.Error())
}

func TestErrGroupGoRoutineFuncError(t *testing.T) {
	g := NewErrGroupWithRecovery(context.Background())
	g.Go(func() error {
		return errors.New("Dummy error")
	})

	err := g.Wait()
	assert.EqualError(t, err, "Dummy error")
}

func TestErrGroupFinishWithoutError(t *testing.T) {
	g := NewErrGroupWithRecovery(context.Background())
	g.Go(func() error {
		time.Sleep(time.Second)
		return nil
	})

	err := g.Wait()
	assert.NoError(t, err)
}

func TestNewErrGroupWithRecoveryAndSharedMutex(t *testing.T) {
	g := NewErrGroupWithRecoveryAndSharedMutex(context.Background())
	g.Go(func(m *sync.Mutex) error {
		m.Lock()
		defer m.Unlock()
		fmt.Println("First line!")
		return nil
	})

	g.Go(func(m *sync.Mutex) error {
		m.Lock()
		defer m.Unlock()
		fmt.Println("Second line!")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
