package mempool

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	KeyHash = "TxHash"
	KeyGas = "Gas"
	KeyFee = "FeePerGas"
	KeySignature = "Signature"

	ErrFieldNotFound = "Field %s not found in line [%s]"
	ErrInvalidToken = "Invalid token [%s]"
	ErrInvalidValueForField = "Invalid value for field %s [%s]"
)

type Transaction struct {
	Hash string
	Gas int
	// The fee is of string type to avoid altering the value when it is being written into the output
	// due to float type not being able to accurately represent decimal fraction
	FeePerGas string
	Signature string

	feePerGasNumeric float64
	totalFee float64
}

func (t Transaction) Priority() float64 {
	return t.totalFee
}

func (t Transaction) Fee() float64 {
	return t.feePerGasNumeric * float64(t.Gas)
}

func (t Transaction) String() string {
	format := fmt.Sprintf("%s=%%s %s=%%v %s=%%s %s=%%s", KeyHash, KeyGas, KeyFee, KeySignature)
	return fmt.Sprintf(format, t.Hash, t.Gas, t.FeePerGas, t.Signature)
}

func ReadTransaction(line string) (Transaction, error) {
	var tx Transaction
	var err error
	tokens := strings.Fields(line)

	tokensMap := make(map[string]string)
	for _, token := range tokens {
		keyVal := strings.Split(token, "=")
		if len(keyVal) != 2 {
			return tx, errors.New(fmt.Sprintf(ErrInvalidToken, token))
		}
		tokensMap[keyVal[0]] = keyVal[1]
	}

	tx, err = transactionFromTokensMap(tokensMap, line)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func transactionFromTokensMap(tokensMap map[string]string, line string) (Transaction, error) {
	var tx Transaction
	var err error
	var ok bool

	tx.Hash, ok = tokensMap[KeyHash]
	if !ok {
		return tx, errors.New(fmt.Sprintf(ErrFieldNotFound, KeyHash, line))
	}

	gasStr, ok := tokensMap[KeyGas]
	if !ok {
		return tx, errors.New(fmt.Sprintf(ErrFieldNotFound, KeyGas, line))
	}

	tx.Gas, err = strconv.Atoi(gasStr)
	if err != nil {
		return tx, errors.New(fmt.Sprintf(ErrInvalidValueForField, KeyGas, gasStr))
	}

	tx.FeePerGas, ok = tokensMap[KeyFee]
	if !ok {
		return tx, errors.New(fmt.Sprintf(ErrFieldNotFound, KeyFee, line))
	}
	if tx.feePerGasNumeric, err = strconv.ParseFloat(tx.FeePerGas, 64); err != nil {
		return tx, errors.New(fmt.Sprintf(ErrInvalidValueForField, KeyFee, tx.FeePerGas))
	}
	tx.totalFee = tx.feePerGasNumeric * float64(tx.Gas)

	tx.Signature, ok = tokensMap[KeySignature]
	if !ok {
		return tx, errors.New(fmt.Sprintf(ErrFieldNotFound, KeySignature, line))
	}

	return tx, nil
}