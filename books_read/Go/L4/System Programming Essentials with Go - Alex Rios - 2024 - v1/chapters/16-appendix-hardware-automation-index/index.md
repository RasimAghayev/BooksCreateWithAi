# Appendix — Hardware Automation (Avadanlıq Avtomatlaşdırması)

## Bu chapter nədən bəhs edir?
Fiziki avadanlıqla işə: USB flash disk hadisələri (D-Bus/UDisks2, /proc/mounts, faylları uzantıya görə avtomatik səlahiyyətləndirmə), sistem bildirişləri (freedesktop.org Notifications), Bluetooth hadisələri (RSSI ilə smartwatch yaxınlığı — ekranda kilidləmə) və XDG/freedesktop.org standartları.

## Əsas fikirlər

### 1. USB — fayl təşkilatçısı
**Ssenari:** flash diskdəki qarışıq kök qovluğunu uzantıya görə avtomatik səlahiyyətləndirən proqram (music_2.wav → wav/, Book_2009.pdf → pdf/).

**Təşkilatlandırıcı funksiya:**
```go
func organizeFiles(paths []string) ([]string, error) {
    events := make([]string, 0)
    for _, path := range paths {
        err := filepath.WalkDir(path, func(path string, dir os.DirEntry, err error) error {
            if err != nil { return err }
            if !dir.IsDir() {
                ext := filepath.Ext(path)
                destDir := filepath.Join(filepath.Dir(path), ext[1:])   // ".pdf" → "pdf"
                destPath := filepath.Join(destDir, dir.Name())
                if err := os.MkdirAll(destDir, os.ModePerm); err != nil { return err }
                if err := os.Rename(path, destPath); err != nil { return err }
                events = append(events, fmt.Sprintf("Moved %s to %s\n", path, destPath))
            }
            return nil
        })
    }
    return events, err
}
```
**Sub-kod izahı:** WalkDir ağacı gəzir; uzantıdan ( nöqtəsiz) qovluq düzəlir; MkdirAll yaradır (varsa xəta vermir); Rename faylı daşıyır; bütün hərəkətlər events jurnalına yazılır.

**Test helper-i:** `os.CreateTemp(dir, "*"+ext)` — pattern-dəki son `*` random soruşulur (`1217776936.txt` kimi).

### 2. Mount nöqtəsini tapmaq — /proc/mounts
Kernel-in real-vaxt virtual faylı (diskdə deyil). Sətir formatı — 6 sahə:
`sahə1 device/UUID | sahə2 mount nöqtəsi | sahə3 fayl sistemi tipi (ext4/ntfs/tmpfs) | sahə4 mount seçimləri (rw/ro) | sahə5 dump bayrağı (0/1) | sahə6 fsck sırası`

**Oxuyan proqram:**
```go
path := os.Args[1]                       // /dev/ ilə başlamalı
if !strings.HasPrefix(path, "/dev/") { return }

file, err := os.Open("/proc/mounts")
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    fields := strings.Fields(scanner.Text())
    device, mountPoint := fields[0], fields[1]
    if strings.HasPrefix(device, path) {
        mountPoint = strings.ReplaceAll(mountPoint, "\\040", " ")  // boşluq kodu!
        filepath.Walk(mountPoint, func(p string, info os.FileInfo, err error) error {
            fmt.Println(p)
            return nil
        })
    }
}
```
**Sub-kod izahı:** `\\040` = `/proc/mounts`-da boşluqun escape kodu ( `/media/My Drive` → `/media/My\040Drive`) — sahə ayırıcı space olduğundan. Walk mount nöqtəsindən bütün faylları siyahılayır.
**Manual üsul:** `df -h` → `/dev/sdc1 15G 16M 15G 1% /media/alexrios/usbtest`.

### 3. Anlayışlar: Partition vs Block vs Device vs Disk
| | Nədir | Nümunə |
|---|---|---|
| **Disk** | fiziki saxlama avadanlığı | HDD, SSD, optical, NAS |
| **Partition** | diskin məntiqi bölgüsü — OS izolyasiyası, data təşkili, təhlükəsizlik | /dev/sdc1 |
| **Block** | sabit ölçülü data vahidi — OS bunlarla oxuyub-yazır | 512B / 4KB |
| **Device** | fiziki VƏYA virtual saxlama — OS-un interfeysi | /dev/sdc |

**/dev/sdc** = bütün fiziki disk (format/partition əməliyyatları); **/dev/sdc1** = onun 1-ci partition-u (mount/fayl sistemi əməliyyatları).

### 4. freedesktop.org və XDG
**freedesktop.org** (2000, Havoc Pennington): GNOME/KDE/layihələr arasında interoperabilitet missiyası. **XDG** (X Desktop Group): desktop-lar arası standartlar — Base Directory, Desktop Menu spec. **D-Bus**: onların ekosistemindən çıxan message bus — tətbiqlərin bir-biri və sistemlə danışığı (Linux-ın poçt xidməti).

### 5. D-Bus — sistem bildirişləri
Kitabxana: `github.com/godbus/dbus/v5`.

```go
conn, err := dbus.ConnectSessionBus()      // session bus (istifadəçi səviyyəsi)
defer conn.Close()

obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
call := obj.Call("org.freedesktop.Notifications.Notify", 0,
    appName, replacesID, appIcon, summary, body, actions, hints, expireTimeout)
```
**Notify parametrləri:**
| Parametr | Tip | Məna |
|---|---|---|
| app_name | STRING | göndərən tətbiqin adı |
| replaces_id | UINT32 | 0 = əvəzləmə yoxdur |
| app_icon | STRING | ikon adı/yolu (boş = yoxdur) |
| summary | STRING | başlıq (mütləq) |
| body | STRING | əsas mətn |
| actions | []string | cüt-cüt düymlər (cüt = ID, tək = mətn) |
| hints | map[string]dbus.Variant | əlavə xassələr (PID, urgency) |
| expire_timeout | INT32 | ms; -1 = server default, 0 = heç vaxt |

### 6. D-Bus — USB hadisələri (UDisks2)
Flash disk taxılanda avtomatik reaksiya üçün system bus + signal match rule.

```go
conn, err := dbus.SystemBus()              // SYSTEM bus — hardware xidmətləri
defer conn.Close()

ch := make(chan *dbus.Signal)             // signal kanalı
conn.Signal(ch)

matchRule := "type='signal',sender='org.freedesktop.UDisks2'," +
    "interface='org.freedesktop.DBus.ObjectManager'," +
    "path='/org/freedesktop/UDisks2'"
conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

for signal := range ch {
    if signal.Name == "org.freedesktop.DBus.ObjectManager.InterfacesAdded" {
        path := signal.Body[0].(dbus.ObjectPath)
        if strings.HasPrefix(string(path), "/org/freedesktop/UDisks2/block_devices/") {
            deviceObj := conn.Object("org.freedesktop.UDisks2", path)
            deviceProps := deviceObj.Call("org.freedesktop.DBus.Properties.Get", 0,
                "org.freedesktop.UDisks2.Block", "Device")
            mountPoints := deviceProps.Body[0].(dbus.Variant)
            fmt.Println(mountPoints.Value())        // /dev/sdc, /dev/sdc1...
        }
    }
}
```
**Sub-kod izahı:** match rule yalnız UDisks2 ObjectManager signal-larını filtrələyir; AddMatch bus daemon-a qaydani qeydiyyatdan keçirir; `InterfacesAdded` = yeni obyekt (block device) əlavə olundu; Properties.Get vasitəsilə cihaz xassələri oxunur.

**Mount nöqtələri D-Bus-dan ( /proc/mounts parse-ləmədən imtina):**
```go
func mountPoints(deviceNames []string) ([]string, error) {
    conn, _ := dbus.ConnectSystemBus()
    defer conn.Close()
    var mountPoints []string
    for _, deviceName := range deviceNames {
        objPath := path.Join("/org/freedesktop/UDisks2/block_devices", deviceName)
        obj := conn.Object("org.freedesktop.UDisks2", dbus.ObjectPath(objPath))
        var result map[string]dbus.Variant
        obj.Call("org.freedesktop.DBus.Properties.GetAll", 0,
            "org.freedesktop.UDisks2.Filesystem").Store(&result)
        if mp, ok := result["MountPoints"]; ok {
            for _, m := range mp.Value().([][]byte) {
                mountPoints = append(mountPoints, string(m))
            }
        }
    }
    return mountPoints, nil
}
```
**Sub-kod izahı:** Filesystem interfeysinin `MountPoints` xassəsi `[][]byte` qaytarır (hər partition-un mount yolu) — kernel virtual faylından oxumaq əvəzinə birbaşa D-Bus property oxunur.

### 7. Bluetooth — smartwatch proximity lock
Kitabxana: `github.com/muka/go-bluetooth/api`. Ssenari: Galaxy Watch Active2 RSSI -70 dBm-dən aşağı düşəndə iş stansiyası kilidlənir.

**Kəşf + polling:**
```go
adapter, err := api.GetDefaultAdapter()
adapter.StartDiscovery()

ticker := time.NewTicker(10 * time.Second)     // 10 saniyəlik polling
defer ticker.Stop()
for {
    select {
    case <-ticker.C:
        devices, _ := adapter.GetDevices()
        for _, device := range devices {
            info, _ := device.GetProperties()
            if info.RSSI < -70 && info.Name == "Galaxy Watch Active2(207D) LE" {
                lockScreen()
            }
        }
    }
}
```
**Sub-kod izahı:** ticker ilə dövri skan (bir dəfəlik yox); şərt: saat tapıldı + siqnal zəif = istifadəçi uzaqda.

**Kilidləmə:**
```go
func lockScreen() error {
    _, err := exec.Command("xdg-screensaver", "lock").Output()
    return err
}
```

### 8. RSSI (Received Signal Strength Indicator)
| dBm | Məna |
|---|---|
| 0 | mənbənin yanında (gossip central) |
| -30 | çox yaxın |
| -70 | təxminən "masadan uzaqlaşma" həddi (calibrasiya ilə tapılır) |
| -100 | praktik olaraq itmiş siqnal |

- **Dəyər:** yaxınlıq təxmini (dəqiq məsafə YOX) — qapı açma, attendance, indoor positioning.
- **Məhdudiyyət:** divar/mikrodalğa/interferensiya dəyişir; moving average filtr yumşaldır; hər cihaz üçün kalibrasiya lazımdır.

### 9. XDG dilemma — Wayland
`xdg-screensaver` XDG standartlaşdırmasının məhsuludur — desktop-lar arası vahid screensaver idarəsi. Amma **Wayland** mühitlərində işləmir — one-size-fits-all həllin həddi. Dərs: desktop avtomatlaşdırmasında portable görünmə həmişə zəmanət deyil.

## Əsas terminlər
- Hardware automation vs software automation — fiziki vs rəqəmsal hadisələr
- filepath.WalkDir / filepath.Ext / os.MkdirAll / os.Rename
- os.CreateTemp(dir, "*"+ext) — pattern-based temp fayl
- /proc/mounts — kernel-in real-vaxt mount virtual faylı (device, mountpoint, fstype, options, dump, fsck)
- `\\040` — /proc/mounts boşluq escape kodu
- df -h — mount tapmağın manual yolu
- Disk / Partition / Block / Device / /dev/sdc vs /dev/sdc1
- freedesktop.org / XDG (X Desktop Group) — desktop interop standartları
- D-Bus: session bus vs system bus
- github.com/godbus/dbus/v5 — ConnectSessionBus / SystemBus / Signal / AddMatch
- Match rule: type/sender/interface/path filtrləri
- org.freedesktop.UDisks2 — disk idarə xidməti
- InterfacesAdded / Properties.Get / Properties.GetAll
- dbus.ObjectPath / dbus.Variant — D-Bus tipləri
- org.freedesktop.Notifications.Notify — sistem bildirişi
- github.com/muka/go-bluetooth/api — GetDefaultAdapter/StartDiscovery/GetDevices
- RSSI — siqnal gücü göstəricisi (0 … -100 dBm)
- xdg-screensaver lock / Wayland uyğunsuzluğu
- time.NewTicker + select — dövri polling (dərs 14-dən fərqli olaraq: burda timer dayandırıla bilir, leak yoxdur)

## Praktik nəticə
1. USB avtomatlaşdırma zənciri: D-Bus signal → mount nöqtəsi (Properties.GetAll) → WalkDir → uzantı qovluqlarına Rename → bildiriş.
2. /proc/mounts parse edərkən `\\040` → boşluq əvəzləməsini unutma — mount path-lər səhv oxunur.
3. /dev/sdc (bütün disk) ilə /dev/sdc1 (partition) fərqini bilmə — format mount səviyyəsində qərarlaır.
4. Hardware hadisələri system bus-dan, istifadəçi UI hadisələri session bus-dan — yanlış bus = siqnal gəlmir.
5. RSSI yaxınlıq ölçüsüdür, məsafə deyil — threshold-u kalibrasiya et, moving average ilə sabitləşdir.
6. Desktop avtomatlaşdırmasında XDG alətləri köhnələ bilər — Wayland kimi yeni mühitlərdə sına.

## Mənbə
Pages: 345-372 (PDF səh. 366-391)
