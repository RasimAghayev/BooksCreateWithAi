# Chapter 13 — Programming from the Command Line (səh. 420-439)

## Bu fəsil nədən bəhs edir?

CLI proqramlaşdırma: os.Args, flag paketi (default + help), stdin/stdout
axını (Rot13 pipeline), exit kodları və best practices, interrupt
siqnalları (signal.Notify), os/exec ilə xarici əmrlər, Bubbletea ilə TUI
və go install.

## Əsas fikirlər

### 1. Arqument oxuma (os.Args)
```go
args := os.Args                 // [0] = proqram adı, [1:] = arqumentlər
if len(args) < 2 {
    fmt.Println("Usage: go run main.go <name>")
    return
}
name := args[1]
greeting := fmt.Sprintf("Hello, %s! Welcome to the command line.", name)
```

### 2. flag paketi
```go
var (
    nameFlag  = flag.String("name", "Sam", "Name of the person to say hello to")
    quietFlag = flag.Bool("quiet", false, "Toggle to be quiet when saying hello")
)

func main() {
    flag.Parse()                    // PARSE — dəyərləri doldurur
    if !*quietFlag {                  // POINTER-dır — * ilə oxu
        greeting := fmt.Sprintf("Hello, %s! Welcome to the command line.", *nameFlag)
        fmt.Println(greeting)
    }
}
```
```bash
go run main.go                          # Hello, Sam! (default-lar)
go run main.go --name=Cassie            # Hello, Cassie!
go run main.go --quiet=true             # (susduruldu)
go run main.go --help                   # AVTOMATİK help:
#  -name string  Name of the person to say hello to (default "Sam")
#  -quiet       Toggle to be quiet when saying hello
```
- flag.String/Bool → default + təsvir; help AVTOMATİK generasiya olunur

### 3. Böyük data axını (stdin/stdout + Rot13)
**Rot13:** hərfi əlifbada 13 irəli sürüşdür (A→N); simmetrik — iki dəfə
tətbiq = orijinal; oyuncaq şifrə (real təhlükəsizlik YOX).

```go
func rot13(s string) string {
    result := make([]byte, len(s))
    for i := 0; i < len(s); i++ {
        char := s[i]
        switch {
        case char >= 'a' && char <= 'z':
            result[i] = 'a' + (char-'a'+13)%26
        case char >= 'A' && char <= 'Z':
            result[i] = 'A' + (char-'A'+13)%26
        default:
            result[i] = char
        }
    }
    return string(result)
}

// Stdin-dən sətir-sətir:
func processStdin() {
    reader := bufio.NewReader(os.Stdin)
    for {
        input, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Error reading stdin:", err)
            return
        }
        fmt.Print(rot13(input))
    }
}

// Fayl YOXDURSA stdin, VARSA fayl:
func processFileOrInput() {
    var inputReader io.Reader
    if len(os.Args) > 1 {
        file, err := os.Open(os.Args[1])
        if err != nil {
            fmt.Println("Error opening file:", err)
            return
        }
        defer file.Close()
        inputReader = file
    } else {
        fmt.Print("Enter text: ")
        inputReader = os.Stdin
    }
    scanner := bufio.NewScanner(inputReader)
    for scanner.Scan() {
        fmt.Println(rot13(scanner.Text()))
    }
}

// ƏSAS: stdin PİP-lənib? (terminal yoxdursa data var!)
stat, _ := os.Stdin.Stat()
if (stat.Mode() & os.ModeCharDevice) == 0 {
    processStdin()      // pip edilmiş data
} else {
    processFileOrInput() // interaktiv
}
```
```bash
echo "enjoy" | go run main.go     # rawbl
cat data.txt | go run main.go     # fayl axını
```
- Fayl oxumadan pip aşkarlanması — ModeCharDevice hiyləsi

### 4. Exit kodları
```go
const (
    ExitCodeSuccess       = 0   // uğur
    ExitCodeInvalidInput  = 1   // xəta
    ExitCodeFileNotFound  = 2   // xəta
)
os.Exit(ExitCodeSuccess)
```
- `echo $?` — son əmrin exit statusu
- Konvensiya: 0 = uğur; qeyri-sıfır = xəta

**Best practices:**
- Ardıcıl logging (mənalı mesajlar)
- Aydın usage məlumatı
- Help və version flag-ləri
- Graceful termination (təmizlik + exit kodu)

### 5. Interrupt siqnalları
**Nədir:** OS proseslərlə siqnal danışır; SIGINT (Ctrl+C), SIGTERM.
Graceful shutdown = resurs buraxma, vəziyyət yadda saxlama,
goroutine-lərin işini bitirməsi.

```go
// signal.Notify kanallara siqnal qeydiyyatı
// (Bu kitabda nümunə göstərilməyib, amma konsept olaraq:
//   ch := make(chan os.Signal, 1)
//   signal.Notify(ch, os.Interrupt)
//   <-ch → cleanup → exit)
```
- Context + timeout ilə shutdown-da donma qarşısı

### 6. Xarici əmrlər (os/exec)
```go
func main() {
    timeLimit := 5 * time.Second
    fmt.Println("Press Enter to start the stopwatch...")
    _, err := fmt.Scanln()               // Enter gözlə
    if err != nil {
        fmt.Println("Error reading from stdin:", err)
        return
    }
    fmt.Println("Stopwatch started. Waiting for", timeLimit)

    time.Sleep(timeLimit)
    fmt.Println("Time's up! Executing the other command.")
    cmd := exec.Command("echo", "Hello")   // əmr + arqumentlər
    cmd.Stdout = os.Stdout                 // çıxışları yönləndir
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        fmt.Println("Error executing command:", err)
    }
}
```
- İşçi qovluq, env dəyişənləri, stdin ötürmə mümkündür
- Cross-platform — OS-dən asılı deyil

### 7. TUI (Bubbletea)
**Konseptlər:** komponentlər (düymə/siyağı), layout, klaviatura hadisələri.

```go
import tea "github.com/charmbracelet/bubbletea"

var choices = []string{"File input", "Type in input"}

type model struct {
    cursor int
    choice string
}

func (m model) Init() tea.Cmd { return nil }

// Update — klaviatura hadisələri:
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q", "esc":
            return m, tea.Quit
        case "enter":
            m.choice = choices[m.cursor]
            return m, tea.Quit
        case "down", "j":                 // vim üslubu!
            m.cursor++
            if m.cursor >= len(choices) {
                m.cursor = 0
            }
        case "up", "k":
            m.cursor--
            if m.cursor < 0 {
                m.cursor = len(choices) - 1
            }
        }
    }
    return m, nil
}

// View — görünüş:
func (m model) View() string {
    s := strings.Builder{}
    s.WriteString("Select if you would like to work with file input or type in input:\n\n")
    for i := 0; i < len(choices); i++ {
        if m.cursor == i {
            s.WriteString("(•) ")
        } else {
            s.WriteString("( ) ")
        }
        s.WriteString(choices[i])
        s.WriteString("\n")
    }
    s.WriteString("\n(press q to quit)\n")
    return s.String()
}

// Main — seçimə görə yönləndir:
p := tea.NewProgram(model{})
m, err := p.Run()
if err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
}
if m, ok := m.(model); ok && m.choice == "File input" {
    processFile("data.txt")
}
if m, ok := m.(model); ok && m.choice == "Type in input" {
    processStdin()
}
```
- Model + Update (mesaj emalı) + View (render) — Elm arxitekturası

### 8. go install
```bash
go install                              # cari layihəni $GOPATH/bin-ə
go install github.com/spf13/cobra-cli@latest   # xarici CLI qur

cobra-cli --help                        # indi QLOBAL əlçatandır
```
- GOOS/GOARCH ilə cross-platform hədəf
- cobra — CLI kitabxana + cobra-cli scaffold generator

## Əsas terminlər
- os.Args — [proqram, arqumentlər...]
- flag.String/Bool + flag.Parse — typed flag-lər + auto help
- Rot13 — 13 sürüşdürmə şifri (simmetrik)
- stdin/stdout streaming — pipeline uyğunluğu
- os.ModeCharDevice — terminal mı / pip mi yoxlaması
- Exit code ($?) — 0 uğur / qeyri-sıfır xəta
- SIGINT / SIGTERM — dayandırma siqnalları
- signal.Notify — kanala siqnal qeydiyyatı
- os/exec.Command — xarici proses
- TUI — terminal user interface
- Bubbletea — Elm-üslubu TUI framework (Model/Update/View)
- go install @latest — qlobal CLI quraşdırma
- cobra / cobra-cli — CLI kitabxana + scaffold

## Praktik nəticə
CLI düzəldəndə: sadə arqument üçün os.Args; konfiqurasiya üçün flag
(default+help birlikdə). Pipeline uyğunluğu: stdin pip-lənibsə (Stat +
ModeCharDevice) avtomatik oxu; Scanner ilə sətir-sətir emal — memory
qənaəti. Exit kodları mütləq (0/1/2...); echo $? ilə yoxla. Ctrl+C üçün
signal.Notify + təmizlik. Xarici əmr: exec.Command + Stdout/Stderr
yönləndirmə. Qarşılıqlı seçimlər üçün Bubbletea TUI; paylanan alətlər
go install @latest.

## Mənbə
Pages: 420-439 (PDF 420-439)
