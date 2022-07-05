/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by Sizegen. DO NOT EDIT.

package evalengine

import (
	"math"
	"reflect"
	"unsafe"

	hack "vitess.io/vitess/go/hack"
)

type cachedObject interface {
	CachedSize(alloc bool) int64
}

func (cached *ArithmeticExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(48)
	}
	// field BinaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.BinaryExpr
	size += cached.BinaryExpr.CachedSize(false)
	// field Op vitess.io/vitess/go/vt/vtgate/evalengine.ArithmeticOp
	if cc, ok := cached.Op.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *BinaryExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(32)
	}
	// field Left vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.Left.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	// field Right vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.Right.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *BindVariable) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field Key string
	size += hack.RuntimeAllocSize(int64(len(cached.Key)))
	return size
}
func (cached *BitwiseExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(48)
	}
	// field BinaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.BinaryExpr
	size += cached.BinaryExpr.CachedSize(false)
	// field Op vitess.io/vitess/go/vt/vtgate/evalengine.BitwiseOp
	if cc, ok := cached.Op.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *BitwiseNotExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	return size
}
func (cached *CallExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(80)
	}
	// field Arguments vitess.io/vitess/go/vt/vtgate/evalengine.TupleExpr
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Arguments)) * int64(16))
		for _, elem := range cached.Arguments {
			if cc, ok := elem.(cachedObject); ok {
				size += cc.CachedSize(true)
			}
		}
	}
	// field Aliases []vitess.io/vitess/go/vt/sqlparser.ColIdent
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Aliases)) * int64(40))
		for _, elem := range cached.Aliases {
			size += elem.CachedSize(false)
		}
	}
	// field Method string
	size += hack.RuntimeAllocSize(int64(len(cached.Method)))
	// field F vitess.io/vitess/go/vt/vtgate/evalengine.builtin
	if cc, ok := cached.F.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *CaseExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(48)
	}
	// field cases []vitess.io/vitess/go/vt/vtgate/evalengine.WhenThen
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.cases)) * int64(32))
		for _, elem := range cached.cases {
			size += elem.CachedSize(false)
		}
	}
	// field Else vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.Else.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *CollateExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	return size
}
func (cached *Column) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	return size
}
func (cached *ComparisonExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(48)
	}
	// field BinaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.BinaryExpr
	size += cached.BinaryExpr.CachedSize(false)
	// field Op vitess.io/vitess/go/vt/vtgate/evalengine.ComparisonOp
	if cc, ok := cached.Op.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *ConvertExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(64)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	// field Type string
	size += hack.RuntimeAllocSize(int64(len(cached.Type)))
	return size
}
func (cached *ConvertUsingExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	return size
}
func (cached *EvalResult) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(96)
	}
	// field expr vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.expr.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	// field env *vitess.io/vitess/go/vt/vtgate/evalengine.ExpressionEnv
	size += cached.env.CachedSize(true)
	// field bytes_ []byte
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.bytes_)))
	}
	// field tuple_ *[]vitess.io/vitess/go/vt/vtgate/evalengine.EvalResult
	if cached.tuple_ != nil {
		size += int64(24)
		size += hack.RuntimeAllocSize(int64(cap(*cached.tuple_)) * int64(88))
		for _, elem := range *cached.tuple_ {
			size += elem.CachedSize(false)
		}
	}
	// field decimal_ vitess.io/vitess/go/vt/vtgate/evalengine/internal/decimal.Decimal
	size += cached.decimal_.CachedSize(false)
	return size
}

//go:nocheckptr
func (cached *ExpressionEnv) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(64)
	}
	// field BindVars map[string]*vitess.io/vitess/go/vt/proto/query.BindVariable
	if cached.BindVars != nil {
		size += int64(48)
		hmap := reflect.ValueOf(cached.BindVars)
		numBuckets := int(math.Pow(2, float64((*(*uint8)(unsafe.Pointer(hmap.Pointer() + uintptr(9)))))))
		numOldBuckets := (*(*uint16)(unsafe.Pointer(hmap.Pointer() + uintptr(10))))
		size += hack.RuntimeAllocSize(int64(numOldBuckets * 208))
		if len(cached.BindVars) > 0 || numBuckets > 1 {
			size += hack.RuntimeAllocSize(int64(numBuckets * 208))
		}
		for k, v := range cached.BindVars {
			size += hack.RuntimeAllocSize(int64(len(k)))
			size += v.CachedSize(true)
		}
	}
	// field Row []vitess.io/vitess/go/sqltypes.Value
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Row)) * int64(32))
		for _, elem := range cached.Row {
			size += elem.CachedSize(false)
		}
	}
	// field Fields []*vitess.io/vitess/go/vt/proto/query.Field
	{
		size += hack.RuntimeAllocSize(int64(cap(cached.Fields)) * int64(8))
		for _, elem := range cached.Fields {
			size += elem.CachedSize(true)
		}
	}
	return size
}

//go:nocheckptr
func (cached *InExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(48)
	}
	// field BinaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.BinaryExpr
	size += cached.BinaryExpr.CachedSize(false)
	// field Hashed map[uintptr]int
	if cached.Hashed != nil {
		size += int64(48)
		hmap := reflect.ValueOf(cached.Hashed)
		numBuckets := int(math.Pow(2, float64((*(*uint8)(unsafe.Pointer(hmap.Pointer() + uintptr(9)))))))
		numOldBuckets := (*(*uint16)(unsafe.Pointer(hmap.Pointer() + uintptr(10))))
		size += hack.RuntimeAllocSize(int64(numOldBuckets * 144))
		if len(cached.Hashed) > 0 || numBuckets > 1 {
			size += hack.RuntimeAllocSize(int64(numBuckets * 144))
		}
	}
	return size
}
func (cached *IsExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(32)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	return size
}
func (cached *LikeExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(64)
	}
	// field BinaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.BinaryExpr
	size += cached.BinaryExpr.CachedSize(false)
	// field Match vitess.io/vitess/go/mysql/collations.WildcardPattern
	if cc, ok := cached.Match.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *Literal) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(96)
	}
	// field Val vitess.io/vitess/go/vt/vtgate/evalengine.EvalResult
	size += cached.Val.CachedSize(false)
	return size
}
func (cached *LogicalExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(64)
	}
	// field BinaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.BinaryExpr
	size += cached.BinaryExpr.CachedSize(false)
	// field opname string
	size += hack.RuntimeAllocSize(int64(len(cached.opname)))
	return size
}
func (cached *NegateExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	return size
}
func (cached *NotExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	// field UnaryExpr vitess.io/vitess/go/vt/vtgate/evalengine.UnaryExpr
	size += cached.UnaryExpr.CachedSize(false)
	return size
}
func (cached *UnaryExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(16)
	}
	// field Inner vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.Inner.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *WeightStringCallExpr) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(48)
	}
	// field String vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.String.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	// field Cast string
	size += hack.RuntimeAllocSize(int64(len(cached.Cast)))
	return size
}
func (cached *WhenThen) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(32)
	}
	// field when vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.when.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	// field then vitess.io/vitess/go/vt/vtgate/evalengine.Expr
	if cc, ok := cached.then.(cachedObject); ok {
		size += cc.CachedSize(true)
	}
	return size
}
func (cached *builtinMultiComparison) CachedSize(alloc bool) int64 {
	if cached == nil {
		return int64(0)
	}
	size := int64(0)
	if alloc {
		size += int64(24)
	}
	// field name string
	size += hack.RuntimeAllocSize(int64(len(cached.name)))
	return size
}
