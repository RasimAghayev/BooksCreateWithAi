# Chapter 16 — Command Line Interface (səh. 349-380)

## Bu fəsil nədən bəhs edir?

Cobra kitabxanası ilə peşəkar CLI qurulması: Command tipi, avtomatik help,
typed flags (String/Int/Bool, Var variantları, required), persistent vs local
flags, subcommand-lar, hooks (PreRun/PostRun), xüsusi help/usage, sənəd
generasiyası (Man/Markdown/ReST/YAML), cobra generator və shell completion.

## Əsas fikirlər

### 1. Cobra əsasları
**Nədir:** Kubernetes, GitHub CLI, Istio tərəfindən istifadə olunan CLI
kitabxanası: avtomatik help, shell autocompletion, POSIX-uyğun flag-lər.

**Quraşdırma:** `go get github.com/spf13/cobra`

**Kitabdan kod nümunəsi (minimal CLI):**
```go
var RootCmd = &cobra.Command{
    Use:     "hello",
    Short:   "short message",
    Long:    "Long message",
    Version: "v0.1.0",
    Example: "this is an example",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Save the world with Go!!!")
    },
}

func main() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**Avtomatik help (`--help`):**
```
Long message
Usage:
hello [flags]
Examples:
this is an example
Flags:
-h, --help      help for hello
-v, --version   version for hello
```
- Konvensiya: hər komanda `cmd/` qovluğunda ayrı faylda; main.go yalnız
  `RootCmd.Execute()` çağırır

### 2. Arqumentlər və flag-lər

**Adi arqumentlər (typsız):**
```go
Run: func(cmd *cobra.Command, args []string) {
    fmt.Printf("%s\n", strings.Join(args, ","))
}
// ./main 1 2 three 4 five → 1,2,three,4,five
```

**Typed flag-lər (init-də qeydiyyat):**
```go
var Msg *string

func init() {
    Msg = RootCmd.Flags().String("msg", "Save the world with Go!!!",
        "Message to show")
}

Run: func(cmd *cobra.Command, args []string) {
    fmt.Printf("[[—%s—]]\n", *Msg)
}
// ./main --msg Hello   → [[—Hello—]]
// ./main               → [[—Save the world with Go!!!—]]
// ./main --message X   → Error: unknown flag: --message (+ help çapı)
```

**Çoxlu typed flag:**
```go
var Msg *string
var Rep *int

func init() {
    Msg = RootCmd.Flags().String("msg", "Save the world with Go!!!", "Message to show")
    Rep = RootCmd.Flags().Int("rep", 1, "Number of times to show the message")
}
// ./main --msg Hello --rep 3 → 3 dəfə çap
```

**Required flag:**
```go
RootCmd.MarkFlagRequired("rep")
// ./main → Error: required flag(s) "rep" not set
```

**Var variantı — Config struct-a birbaşa:**
```go
type Config struct {
    Msg string
    Rep int
}
var cnfg Config = Config{}

func init() {
    RootCmd.Flags().StringVar(&cnfg.Msg, "msg", "Save the world with Go!!!", "Message to show")
    RootCmd.Flags().IntVar(&cnfg.Rep, "rep", 1, "Number of times to show the message")
    RootCmd.MarkFlagRequired("rep")
}
// Run içində: cnfg.Rep, cnfg.Msg — pointer dereference lazım deyil
```

### 3. Komandalar (subcommand)
**Nədir:** Root komanda özü işləməyən "konteyner" olur; əməliyyatlar
subcommand-lardadır.

**Kitabdan kod nümunəsi:**
```go
var RootCmd = &cobra.Command{
    Use:  "say",
    Long: "Root command",
}

var HelloCmd = &cobra.Command{
    Use:   "hello",
    Short: "Say hello",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello!!!")
    },
}

var ByeCmd = &cobra.Command{
    Use:   "bye",
    Short: "Say goodbye",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Bye!!!")
    },
}

func init() {
    RootCmd.AddCommand(HelloCmd, ByeCmd)
}
```
- `./say` → help göstərir (Available Commands: bye, hello, help)
- `./say hello` → Hello!!!; `./say bye --help` → bye-nin öz help-i

### 4. Persistent vs local flags
**PersistentFlags** → uşaqlara miras olunur (qlobal); **Flags** → yalnız
lokal.

```go
var msg string
var person string

func init() {
    RootCmd.AddCommand(HelloCmd, ByeCmd, CustomCmd)

    RootCmd.PersistentFlags().StringVar(&person, "person", "Mr X", "Receiver")
    CustomCmd.Flags().StringVar(&msg, "msg", "what's up", "Custom message")
}
// ./say bye --person John        → Bye John!!!  (persistent işləyir)
// ./say custom --person John     → Say what's up to John
// ./say custom --help → "Global Flags: --person" bölməsi görünür
```

### 5. Hooklar (PreRun/PostRun)
**Nədir:** Run-dan əvvəl/sonra icra olunan funksiyalar.
`PersistentPreRun`/`PersistentPostRun` uşaqlara da miras olur.

```go
var RootCmd = &cobra.Command{
    Use: "say",
    Long: "Root command",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Hello %s!!!\n", person)
    },
    Run: func(cmd *cobra.Command, args []string) {},
    PostRun: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Bye %s!!!\n", person)
    },
}

var SomethingCmd = &cobra.Command{
    Use: "something",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("%s\n", msg)
    },
    PostRun: func(cmd *cobra.Command, args []string) {   // root-un PostRun-u əvəz edir
        fmt.Printf("That's all I have to say %s\n", person)
    },
}
// ./say                          → Hello Mr X!!! / Bye Mr X!!!
// ./say something --msg "How are you?"
//   → Hello Mr X!!! (persistent pre) / How are you? / That's all... (öz post)
```

### 6. Əlavə imkanlar

**Xüsusi help/usage:**
```go
func helper(cmd *cobra.Command, args []string) {
    fmt.Printf("You entered command %s\n", cmd.Name())
    fmt.Println("And that is all the help we have right now :)")
}

func usager(cmd *cobra.Command) error {
    fmt.Printf("You entered command %s\n", cmd.Name())
    fmt.Println("And you do not know how it works :)")
    return errors.New("Something went wrong :(")
}

RootCmd.SetHelpFunc(helper)
RootCmd.SetUsageFunc(usager)
```
- `Args: cobra.MinimumNArgs(2)` → minimum 2 arqument tələbi

**Sənəd generasiyası (cobra/doc):**
```go
import "github.com/spf13/cobra/doc"

header := &doc.GenManHeader{Title: "Test", Manual: "MyManual", Section: "1"}

doc.GenManTree(RootCmd, header, ".")     // man pages
doc.GenMarkdownTree(RootCmd, ".")        // markdown
doc.GenReSTTree(RootCmd, ".")            // reST
doc.GenYamlTree(RootCmd, ".")            // YAML
```
Nümunə YAML çıxışı:
```yaml
name: test
synopsis: Documented test
description: How to document a command
usage: test [flags]
options:
- name: flag
  default_value: "true"
  usage: Some flag
example: ./main test
```

**Cobra generator (skeleton yaratma):**
```bash
>>> $GOPATH/bin/cobra init --pkg-name github.com/.../example_03
>>> $GOPATH/bin/cobra add test
# main.go + cmd/ şablonları + license avtomatik yaranır
```

**Shell completion (Bash/Zsh/Fish/PowerShell):**
```go
var CompletionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate completion script",
    DisableFlagsInUseLine: true,
    ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
    Args: cobra.ExactValidArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        switch args[0] {
        case "bash":
            cmd.Root().GenBashCompletion(os.Stdout)
        case "zsh":
            cmd.Root().GenZshCompletion(os.Stdout)
        case "fish":
            cmd.Root().GenFishCompletion(os.Stdout, true)
        case "powershell":
            cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
        }
    },
}
```
```bash
>>> ./say completion bash > /tmp/completion
>>> source /tmp/completion
>>> ./say [tab][tab]      # bye / completion / hello / help siyahısı
```

**Dinamik ValidArgsFunction (runtime arqumentlər — məs. DB-dən):**
```go
var GetCmd = &cobra.Command{
    Use: "get",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Get user %s!!!\n", args[0])
    },
    ValidArgsFunction: UserGet,   // runtime funksiya
}

func UserGet(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    rand.Seed(time.Now().UnixNano())
    if rand.Int()%2 == 0 {
        return []string{"John", "Mary"}, cobra.ShellCompDirectiveNoFileComp
    }
    return []string{"Ernest", "Rick", "Mary"}, cobra.ShellCompDirectiveNoFileComp
}
// ./db get [tab][tab] → John Mary (və ya Ernest Rick Mary)
```
- `ShellCompDirective` → shell davranışını dəyişən binary flag

## Əsas terminlər
- Cobra — Go CLI framework
- Command — CLI əmri (Use/Short/Run sahələri ilə)
- Flag — komanda parametri (--msg kimi)
- Persistent flag — uşaq komandalara miras olan flag
- Hook (PreRun/PostRun) — Run ətrafı icra qatları
- ValidArgs/ValidArgsFunction — etibarlı arqument siyahısı (statik/dinamik)
- Shell completion — tab-tamamlama skripti
- cobra init/add — skeleton generatoru

## Praktik nəticə
Cobra ilə CLI = Command-lar ağacı. Root boş konteyner, subcommand-lar
əməliyyatlardır. Flag-ləri init-də qeyd edin: `Flags().StringVar(&cfg.X, ...)`
üslubu Config struct ilə ən təmizidir; məcburi flag-lər üçün
`MarkFlagRequired`. Qlobal parametrlər `PersistentFlags`-ə. Help avtomatik;
lazımsa SetHelpFunc ilə dəyişin. `cobra/doc` ilə sənədləri (man/md/yaml)
generasiya edin; `cobra init/add` ilə başlayın; completion komandası ilə
istifadəçi təcrübəsini artırın. Dinamik arqumentlər (DB istifadəçiləri kimi)
`ValidArgsFunction` ilə verilir.

## Mənbə
Pages: 349-380 (PDF 349-380)
