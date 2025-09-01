/*
 * Copyright 2025 Pouya Vedadiyan
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
 */
package exql

import (
	"github.com/vedadiyan/exql/lang"
)

func Parse(expr string) (lang.ExprNode, error) {

	return lang.ParseExpression(expr)
}

func Eval(expr string, context lang.Context) (lang.Value, error) {
	result, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	return result.Evaluate(context)
}
