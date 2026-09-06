# Get Programming with Go — Terminologiya lüğəti (AZ)

`English (Azərbaycanca)` formatında — SYSTEM_PROMPT qayda 4 üzrə.

## Unit 0-1: Başlanğıc + İmperativ
- Compiler / Interpreter (kompilyator / tərcüməçi)
- Package / import / func
- One True Brace Style (tək doğru brace üslubu)
- Format Verb (`%v`, `%T`, `%c`, `%f`, `%b`, `%x`)
- Escape Sequence (`\n`)
- Padding / Precision / Width (`%05.2f`)
- Constant / Variable (const/var)
- Magic Number (maqik rəqəm)
- Modulus Operator (`%`)
- Off-by-one Error (bir-vahid xətası)
- Pseudorandom (psevdotəsadüfi — rand.Intn)
- Boolean / Comparison Operator
- Short-circuit Logic (&& / || erkən dayanma)
- Fallthrough (switch-də eksplisit keçid)
- Infinite Loop (`for {}`)
- Variable Scope (package/function/block)
- Short Declaration (`:=`)
- Code Smell / Refactoring

## Unit 2: Tiplər
- IEEE-754 Floating-point (üzən nöqtə)
- float64 / float32 (double / single precision)
- Zero Value (sıfır dəyər — tipin default-u)
- Rounding Error (yuvarlama xətası)
- Machine Epsilon (2⁻⁵²)
- Signed / Unsigned (işarəli/işarəsiz)
- Wrap-around (dövrü taşma — 255+1=0)
- Hexadecimal (`0x`) / Binary (`%b`) / Nibble
- Exponent Notation (`41.3e12`)
- big.Int / big.Float / big.Rat
- Untyped Constant / Literal (tipsiz sabit)
- Raw String Literal (backtick — `\` işləməyən)
- Unicode / Code Point / UTF-8
- rune (int32 alias) / byte (uint8 alias)
- Character Literal ('A')
- len (bayt!) vs RuneCountInString (simvol)
- Blank Identifier (`_`)
- Type Conversion (eksplisit; koercion YOX)
- Static vs Dynamic Typing
- strconv: Itoa / Atoi

## Unit 3: Building Blocks
- Parameter vs Argument
- Exported (böyük hərf) / Unexported
- Multiple Return Values / Named Result
- Variadic Function (`...`)
- Empty Interface (`interface{}`)
- Pass by Value (kopya ilə ötürmə)
- Side-effect-free Function (təsirsiz funksiya)
- New Type (alias-dan fərqli: `type celsius float64`)
- Method / Receiver (dot notation)
- First-class Function
- Function Signature / Function Type
- Anonymous Function / Function Literal
- Closure (əhatə — dəyişənlərə referans)

## Unit 4: Kolleksiyalar
- Array (sabit uzunluq = tipin hissəsi)
- Composite Literal
- Bounds / Panic (`index out of range`)
- Slice (array-ə pəncərə)
- Half-open Range (`[0:4)`)
- Default Indices (`s[:4]`, `s[:]`)
- Underlying Array (altta yatan massiv)
- Shared View (paylaşılan görünüş — dəyişiklik görünür)
- append / len / cap
- Three-index Slicing (`s[i:j:k]` — cap məhdudlaşdırma)
- make (preallokasiya)
- Ellipsis 3 istifadəsi: `[...]`, `...T`, `slice...`
- Map (Python dict, Ruby hash, JS object)
- Comma-ok Idiom (`v, ok := m[k]`)
- delete (builtin)
- Set İmprovizasiyası (map[T]bool)
- Double Buffering (Step(a, b) + swap)

## Unit 5: State və Behavior
- Structure / Field
- `%+v` / `%#v` (sahə adları / Go representasiyası)
- Struct Tag (`json:"..."` — raw string)
- json.Marshal / MarshalIndent
- Constructor Function (newType/NewType konvensiyası)
- Composition (HAS-A)
- Embedding (sahə adsız — avtomatik forwarding)
- Method Forwarding / Promotion
- Ambiguous Selector (eyni adlı metod toqquşması)
- Inheritance (IS-A) vs Composition
- Interface / Method Set
- Implicit Satisfaction (implements YOX)
- Polymorphism ("many shapes")
- `-er` Suffix (talker, Stringer)
- fmt.Stringer / io.Reader / io.Writer / json.Marshaler
- Interface Embedding (io.ReadWriter)

## Unit 6: Gopher Hole
- Address Operator (`&`) / Dereference (`*`)
- Pointer Type (`*int`)
- Pointer Arithmetic (Go-da QADAĞAN)
- Automatic Dereference (struct sahə / array)
- Interior Pointer (`&player.stats`)
- Pointer Receiver (tutum qaydası — hamısı pointer!)
- Immutable Type (time.Time)
- nil (pointer/slice/map/interface zero value)
- Nil Pointer Dereference (panic)
- Guard Clause (`if p == nil { return }`)
- "Billion Dollar Mistake" (Tony Hoare)
- Nil Slice (append/range OK) / Nil Map (oxu OK, yazı panic)
- Interface = type + value (`(*int)(nil)` tələsi)
- error Interface / Custom Error Type
- Sentinel Error (`ErrBounds`)
- Multiple Errors (SudokuError — []error)
- Type Assertion (`err.(SudokuError)`)
- defer (funksiya qayıtmazdan əvvəl zəmanətli icra)
- safeWriter Pattern ("errors are values" — Rob Pike)
- panic / recover (yalnız defer daxilində)
- Go Proverbs

## Unit 7: Concurrency
- Goroutine
- Time Sharing (məhdud CPU bölüşdürməsi)
- Channel (`make(chan T)`)
- Send / Receive (`c <- v` / `v = <-c`)
- Blocking (resurssuz gözləmə)
- Deadlock (`<-c` — göndərən yoxdur)
- Nil Channel (əbədi blok; close → panic)
- select (çox kanala baxış)
- time.After (timeout/periodik kanal)
- Timeout Pattern
- Sentinel Value (bitiş siqnalı — "" kimi)
- close / Two-value Receive (`v, ok`)
- range Channel (bağlanana qədər iterasiya)
- Pipeline (source → filter → print)
- Race Condition / Undefined Behavior
- Race Detector (`go run -race`)
- Mutex (mutual exclusion) / Lock / Unlock
- Long-lived Worker (`for { select {} }`)
- Command Channel (API arxasında gizli kanal)
- Event Loop (Node.js) vs Goroutine-per-worker
- Message-passing Architecture (Curiosity rover modulları)
