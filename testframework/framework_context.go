/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package testframework

import (
	"bytes"
	"encoding/hex"
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"github.com/ontio/ontology-crypto/keypair"
	s "github.com/ontio/ontology-crypto/signature"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-test/common"
	"github.com/ontio/ontology/account"
	"math/big"
)

//TestFrameworkContext is the context for test case
type TestFrameworkContext struct {
	Ont       *sdk.OntologySdk //sdk to ontology
	Wallet    account.Client   // wallet instance
	failNowCh chan interface{}
}

//NewTestFrameworkContext return a TestFrameworkContext instance
func NewTestFrameworkContext(ont *sdk.OntologySdk, wal account.Client, failNowCh chan interface{}) *TestFrameworkContext {
	return &TestFrameworkContext{
		Ont:       ont,
		Wallet:    wal,
		failNowCh: failNowCh,
	}
}

//LogInfo log info in test case
func (this *TestFrameworkContext) LogInfo(arg0 interface{}, args ...interface{}) {
	log4.Info(arg0, args...)
}

//LogError log error info  when error occur in test case
func (this *TestFrameworkContext) LogError(arg0 interface{}, args ...interface{}) {
	log4.Error(arg0, args...)
}

//LogWarn log warning info in test case
func (this *TestFrameworkContext) LogWarn(arg0 interface{}, args ...interface{}) {
	log4.Warn(arg0, args...)
}

func (this *TestFrameworkContext) GetDefaultAccount() (*account.Account, error) {
	return this.Wallet.GetDefaultAccount([]byte(common.DefConfig.Password))
}

func (this *TestFrameworkContext) GetAccount(addr string) (*account.Account, error) {
	acc, err := this.Wallet.GetAccountByAddress(addr, []byte(common.DefConfig.Password))
	if err != nil {
		return nil, err
	}
	if acc != nil {
		return acc, nil
	}
	return this.Wallet.GetAccountByLabel(addr, []byte(common.DefConfig.Password))
}

func (this *TestFrameworkContext) NewAccount(label ...string) (*account.Account, error) {
	label_tag := ""
	if len(label) > 0 {
		label_tag = label[0]
	}
	return this.Wallet.NewAccount(label_tag, keypair.PK_ECDSA, keypair.P256, s.SHA256withECDSA, []byte(common.DefConfig.Password))
}

//FailNow will stop test, and skip all haven't not test case
func (this *TestFrameworkContext) FailNow() {
	select {
	case <-this.failNowCh:
	default:
		close(this.failNowCh)
	}
}

func (this *TestFrameworkContext) GetGasPrice() uint64 {
	return common.DefConfig.GasPrice
}

func (this *TestFrameworkContext) GetGasLimit() uint64 {
	return common.DefConfig.GasLimit
}

//AssertToInt compare with int, if not equal, return error
func (this *TestFrameworkContext) AssertToInt(value interface{}, expect int) error {
	v, ok := value.(*big.Int)
	if !ok {
		return fmt.Errorf("Assert:%v to big.Int failed", value)
	}
	if int(v.Int64()) != expect {
		return fmt.Errorf("%v not equal:%v", value, expect)
	}
	return nil
}

//AssertToInt compare with uint, if not equal, return error
func (this *TestFrameworkContext) AssertToUint(value interface{}, expect uint) error {
	v, ok := value.(*big.Int)
	if !ok {
		return fmt.Errorf("Assert:%v to uint failed", value)
	}
	if uint(v.Uint64()) != expect {
		return fmt.Errorf("%v not equal:%v", value, expect)
	}
	return nil
}

//AssertToInt compare with bool, if not equal, return error
func (this *TestFrameworkContext) AssertToBoolean(value interface{}, expect bool) error {
	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("Assert:%v to boolean failed", value)
	}
	if v != expect {
		return fmt.Errorf("%v not equal:%v", value, expect)
	}
	return nil
}

//AssertToInt compare with string, if not equal, return error
func (this *TestFrameworkContext) AssertToString(value interface{}, expect string) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("Assert:%v to string failed", value)
	}
	if v != expect {
		return fmt.Errorf("%v not equal:%v", value, expect)
	}
	return nil
}

//AssertToInt compare with byteArray, if not equal, return error
func (this *TestFrameworkContext) AssertToByteArray(value interface{}, expect []byte) error {
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Assert:%v to string failed", value)
	}
	if !bytes.EqualFold(v, expect) {
		return fmt.Errorf("%x not equal:%x", v, expect)
	}
	return nil
}

//AssertToInt compare with big.Int, if not equal, return error
func (this *TestFrameworkContext) AssertBigInteger(value interface{}, expect *big.Int) error {
	v, ok := value.(*big.Int)
	if !ok {
		return fmt.Errorf("Assert:%v to big.int failed", value)
	}
	if v.Cmp(expect) != 0 {
		return fmt.Errorf("%v not equal:%v", v, expect)
	}
	return nil
}

//ConvertToHexString return hex string
func (this *TestFrameworkContext) ConvertToHexString(v interface{}) (string, error) {
	value, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%v ConvertToString failed", v)
	}
	data, _ := hex.DecodeString(value)
	return string(data), nil
}

//ConvertToHexString return big.Int
func (this *TestFrameworkContext) ConvertToBigInt(v interface{}) (*big.Int, error) {
	value, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("%ConvertToBigInt failed", v)
	}
	data, _ := hex.DecodeString(value)
	return new(big.Int).SetBytes(data), nil
}
