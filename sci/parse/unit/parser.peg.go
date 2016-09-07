package unit

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulefirst
	rulee1
	ruleunits
	ruleinverted_units
	ruleunit
	ruleparens_exp
	ruleparens
	rulediv
	rulemul
	ruleref_exp
	ruleref
	ruleexp
	ruleminus
	ruleopen
	ruleclose
	rulesp
	ruleAction0
	ruleAction1
	rulePegText
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"first",
	"e1",
	"units",
	"inverted_units",
	"unit",
	"parens_exp",
	"parens",
	"div",
	"mul",
	"ref_exp",
	"ref",
	"exp",
	"minus",
	"open",
	"close",
	"sp",
	"Action0",
	"Action1",
	"PegText",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next uint32, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/*func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2 * len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}*/

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type Parser struct {
	Result U
	Err    error
	stack  []U

	Buffer string
	buffer []rune
	rules  [25]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Parser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *Parser) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *Parser) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *Parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.finish()
		case ruleAction1:
			p.invertUnit()
		case ruleAction2:
			p.expParens(buffer[begin:end])
		case ruleAction3:
			p.div()
		case ruleAction4:
			p.mul()
		case ruleAction5:
			p.expUnit(buffer[begin:end])
		case ruleAction6:
			p.ref(buffer[begin:end])

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *Parser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	var tree tokenTree = &tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 first <- <(sp e1 !. Action0)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rulesp]() {
					goto l0
				}
				{
					position2 := position
					depth++
					{
						position3, tokenIndex3, depth3 := position, tokenIndex, depth
						{
							position5 := position
							depth++
							{
								position6, tokenIndex6, depth6 := position, tokenIndex, depth
								if buffer[position] != rune('1') {
									goto l6
								}
								position++
								goto l7
							l6:
								position, tokenIndex, depth = position6, tokenIndex6, depth6
							}
						l7:
							if buffer[position] != rune('/') {
								goto l4
							}
							position++
							if !_rules[rulesp]() {
								goto l4
							}
							if !_rules[ruleunit]() {
								goto l4
							}
						l8:
							{
								position9, tokenIndex9, depth9 := position, tokenIndex, depth
								if !_rules[ruleunit]() {
									goto l9
								}
								goto l8
							l9:
								position, tokenIndex, depth = position9, tokenIndex9, depth9
							}
							{
								add(ruleAction1, position)
							}
							depth--
							add(ruleinverted_units, position5)
						}
						goto l3
					l4:
						position, tokenIndex, depth = position3, tokenIndex3, depth3
						{
							position11 := position
							depth++
							if !_rules[ruleunit]() {
								goto l0
							}
						l12:
							{
								position13, tokenIndex13, depth13 := position, tokenIndex, depth
								if !_rules[ruleunit]() {
									goto l13
								}
								goto l12
							l13:
								position, tokenIndex, depth = position13, tokenIndex13, depth13
							}
							depth--
							add(ruleunits, position11)
						}
					}
				l3:
					depth--
					add(rulee1, position2)
				}
				{
					position14, tokenIndex14, depth14 := position, tokenIndex, depth
					if !matchDot() {
						goto l14
					}
					goto l0
				l14:
					position, tokenIndex, depth = position14, tokenIndex14, depth14
				}
				{
					add(ruleAction0, position)
				}
				depth--
				add(rulefirst, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 e1 <- <(inverted_units / units)> */
		nil,
		/* 2 units <- <unit+> */
		nil,
		/* 3 inverted_units <- <('1'? '/' sp unit+ Action1)> */
		nil,
		/* 4 unit <- <(parens_exp / ref_exp / ((&('*') mul) | (&('/') div) | (&('(') parens) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') ref)))> */
		func() bool {
			position19, tokenIndex19, depth19 := position, tokenIndex, depth
			{
				position20 := position
				depth++
				{
					position21, tokenIndex21, depth21 := position, tokenIndex, depth
					{
						position23 := position
						depth++
						if !_rules[ruleparens]() {
							goto l22
						}
						if buffer[position] != rune('^') {
							goto l22
						}
						position++
						{
							position24 := position
							depth++
							if !_rules[ruleexp]() {
								goto l22
							}
							depth--
							add(rulePegText, position24)
						}
						{
							add(ruleAction2, position)
						}
						depth--
						add(ruleparens_exp, position23)
					}
					goto l21
				l22:
					position, tokenIndex, depth = position21, tokenIndex21, depth21
					{
						position27 := position
						depth++
						{
							position28 := position
							depth++
							{
								position31, tokenIndex31, depth31 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l32
								}
								position++
								goto l31
							l32:
								position, tokenIndex, depth = position31, tokenIndex31, depth31
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l26
								}
								position++
							}
						l31:
						l29:
							{
								position30, tokenIndex30, depth30 := position, tokenIndex, depth
								{
									position33, tokenIndex33, depth33 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l34
									}
									position++
									goto l33
								l34:
									position, tokenIndex, depth = position33, tokenIndex33, depth33
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l30
									}
									position++
								}
							l33:
								goto l29
							l30:
								position, tokenIndex, depth = position30, tokenIndex30, depth30
							}
							if buffer[position] != rune('^') {
								goto l26
							}
							position++
							if !_rules[ruleexp]() {
								goto l26
							}
							depth--
							add(rulePegText, position28)
						}
						if !_rules[rulesp]() {
							goto l26
						}
						{
							add(ruleAction5, position)
						}
						depth--
						add(ruleref_exp, position27)
					}
					goto l21
				l26:
					position, tokenIndex, depth = position21, tokenIndex21, depth21
					{
						switch buffer[position] {
						case '*':
							{
								position37 := position
								depth++
								if buffer[position] != rune('*') {
									goto l19
								}
								position++
								if !_rules[rulesp]() {
									goto l19
								}
								if !_rules[ruleunit]() {
									goto l19
								}
								{
									add(ruleAction4, position)
								}
								depth--
								add(rulemul, position37)
							}
							break
						case '/':
							{
								position39 := position
								depth++
								if buffer[position] != rune('/') {
									goto l19
								}
								position++
								if !_rules[rulesp]() {
									goto l19
								}
								if !_rules[ruleunit]() {
									goto l19
								}
								{
									add(ruleAction3, position)
								}
								depth--
								add(rulediv, position39)
							}
							break
						case '(':
							if !_rules[ruleparens]() {
								goto l19
							}
							break
						default:
							{
								position41 := position
								depth++
								{
									position42 := position
									depth++
									{
										position45, tokenIndex45, depth45 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l46
										}
										position++
										goto l45
									l46:
										position, tokenIndex, depth = position45, tokenIndex45, depth45
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l19
										}
										position++
									}
								l45:
								l43:
									{
										position44, tokenIndex44, depth44 := position, tokenIndex, depth
										{
											position47, tokenIndex47, depth47 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l48
											}
											position++
											goto l47
										l48:
											position, tokenIndex, depth = position47, tokenIndex47, depth47
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l44
											}
											position++
										}
									l47:
										goto l43
									l44:
										position, tokenIndex, depth = position44, tokenIndex44, depth44
									}
									depth--
									add(rulePegText, position42)
								}
								if !_rules[rulesp]() {
									goto l19
								}
								{
									add(ruleAction6, position)
								}
								depth--
								add(ruleref, position41)
							}
							break
						}
					}

				}
			l21:
				depth--
				add(ruleunit, position20)
			}
			return true
		l19:
			position, tokenIndex, depth = position19, tokenIndex19, depth19
			return false
		},
		/* 5 parens_exp <- <(parens '^' <exp> Action2)> */
		nil,
		/* 6 parens <- <(open unit+ close)> */
		func() bool {
			position51, tokenIndex51, depth51 := position, tokenIndex, depth
			{
				position52 := position
				depth++
				{
					position53 := position
					depth++
					if buffer[position] != rune('(') {
						goto l51
					}
					position++
					if !_rules[rulesp]() {
						goto l51
					}
					depth--
					add(ruleopen, position53)
				}
				if !_rules[ruleunit]() {
					goto l51
				}
			l54:
				{
					position55, tokenIndex55, depth55 := position, tokenIndex, depth
					if !_rules[ruleunit]() {
						goto l55
					}
					goto l54
				l55:
					position, tokenIndex, depth = position55, tokenIndex55, depth55
				}
				{
					position56 := position
					depth++
					if buffer[position] != rune(')') {
						goto l51
					}
					position++
					if !_rules[rulesp]() {
						goto l51
					}
					depth--
					add(ruleclose, position56)
				}
				depth--
				add(ruleparens, position52)
			}
			return true
		l51:
			position, tokenIndex, depth = position51, tokenIndex51, depth51
			return false
		},
		/* 7 div <- <('/' sp unit Action3)> */
		nil,
		/* 8 mul <- <('*' sp unit Action4)> */
		nil,
		/* 9 ref_exp <- <(<(([a-z] / [A-Z])+ '^' exp)> sp Action5)> */
		nil,
		/* 10 ref <- <(<([a-z] / [A-Z])+> sp Action6)> */
		nil,
		/* 11 exp <- <(minus? [1-9] [0-9]*)> */
		func() bool {
			position61, tokenIndex61, depth61 := position, tokenIndex, depth
			{
				position62 := position
				depth++
				{
					position63, tokenIndex63, depth63 := position, tokenIndex, depth
					{
						position65 := position
						depth++
						if buffer[position] != rune('-') {
							goto l63
						}
						position++
						depth--
						add(ruleminus, position65)
					}
					goto l64
				l63:
					position, tokenIndex, depth = position63, tokenIndex63, depth63
				}
			l64:
				if c := buffer[position]; c < rune('1') || c > rune('9') {
					goto l61
				}
				position++
			l66:
				{
					position67, tokenIndex67, depth67 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l67
					}
					position++
					goto l66
				l67:
					position, tokenIndex, depth = position67, tokenIndex67, depth67
				}
				depth--
				add(ruleexp, position62)
			}
			return true
		l61:
			position, tokenIndex, depth = position61, tokenIndex61, depth61
			return false
		},
		/* 12 minus <- <'-'> */
		nil,
		/* 13 open <- <('(' sp)> */
		nil,
		/* 14 close <- <(')' sp)> */
		nil,
		/* 15 sp <- <(' ' / '\t')*> */
		func() bool {
			{
				position72 := position
				depth++
			l73:
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					{
						position75, tokenIndex75, depth75 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l76
						}
						position++
						goto l75
					l76:
						position, tokenIndex, depth = position75, tokenIndex75, depth75
						if buffer[position] != rune('\t') {
							goto l74
						}
						position++
					}
				l75:
					goto l73
				l74:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
				}
				depth--
				add(rulesp, position72)
			}
			return true
		},
		/* 17 Action0 <- <{ p.finish() }> */
		nil,
		/* 18 Action1 <- <{ p.invertUnit() }> */
		nil,
		nil,
		/* 20 Action2 <- <{ p.expParens(buffer[begin:end]) }> */
		nil,
		/* 21 Action3 <- <{ p.div() }> */
		nil,
		/* 22 Action4 <- <{ p.mul() }> */
		nil,
		/* 23 Action5 <- <{ p.expUnit(buffer[begin:end]) }> */
		nil,
		/* 24 Action6 <- <{  p.ref(buffer[begin:end]) }> */
		nil,
	}
	p.rules = _rules
}
