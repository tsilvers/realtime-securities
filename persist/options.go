package persist

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/tsilvers/realtime-securities/markets/option"
	"io"
	"os"
)

func InitExpirationsStore() {
	_ = os.RemoveAll(optionsDir)
	_ = os.MkdirAll(optionsDir, 0755)
}

func SaveExpirations(symbol string, exps []option.ExpirationGob) error {
	fname := optionsDir + stockFilename(symbol)
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("could not persist option expirations for %s to file %s: %w", symbol, fname, err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)

	for _, exp := range exps {
		if err := enc.Encode(exp); err != nil {
			return fmt.Errorf("could not persist option expirations for %s to file %s: %w", symbol, fname, err)
		}
	}

	return nil
}

func LoadExpirations(symbol string) ([]*option.ExpirationGob, error) {
	var exps []*option.ExpirationGob

	fname := optionsDir + stockFilename(symbol)
	file, err := os.Open(fname)
	if err != nil {
		return exps, fmt.Errorf("could not load option expirations for %s from file %s: %w", symbol, fname, err)
	}
	defer file.Close()

	enc := gob.NewDecoder(file)

	for err == nil {
		exp := &option.ExpirationGob{}
		if err = enc.Decode(exp); err == nil {
			exps = append(exps, exp)
		}
	}

	if !errors.Is(io.EOF, err) {
		return nil, fmt.Errorf("could not load option expirations for %s from file %s: %w", symbol, fname, err)
	}

	return exps, nil
}
