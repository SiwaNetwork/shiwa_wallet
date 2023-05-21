# Shiwa Wallet (настольный кошелек)

Shiwa Wallet - это официальный кошелек $SHI для сети ShiwaNetwork. Он позволяет пользователям взаимодействовать с сетью различными способами, включая проведение транзакций с $SHI. Кошелек в настоящее время находится в разработке, и уже выпущена первая публичная версия с большими обновлениями в будущем.

### Технологии

Shiwa Wallet использует ряд проектов с открытым исходным кодом для своей работы:

* [VueJS](https://vuejs.org) - Прогрессивный фреймворк JavaScript.
* [Go](https://golang.org) - Go - язык программирования с открытым исходным кодом, который облегчает создание простого, надежного и эффективного программного обеспечения.
* [GORM](https://gorm.io) - Объектно-реляционное отображение для Go.
* [Wails](https://wails.app/) - Фреймворк для создания настольных приложений с использованием Go и веб-технологий.


### Установка

#### Предварительные требования
Для работы Shiwa Wallet требуется [OpenJDK v9](https://java.com/) и [JRE](https://www.oracle.com/java/technologies/javase-jre8-downloads.html). Если вы устанавливаете версию для macOS, используйте [Unarchiver](https://theunarchiver.com/), чтобы распаковать shiwa_installer.zip.

#### Загрузка Shiwa Wallet
Последние версии можно найти в разделе [релизы](https://github.com/.


### Среда разработки

#### 1. Скачайте дистрибутив Go с официального сайта.

Распространение Go и инструменты доступны в виде установщика для всех популярных операционных систем. Посетите <https://golang.org/dl/>, чтобы скачать подходящую версию для вашей операционной системы. Инструкции по установке можно найти [здесь](https://golang.org/doc/install).

#### 2. Скачайте и установите NPM.

NPM и Node.js можно загрузить с официального сайта [здесь](https://nodejs.org/en/download/). Просто выберите вашу дистрибуцию/ОС и архитектуру ЦП.

#### 3. Установите Wails.

Shiwa Wallet создан

 с использованием легковесного фреймворка для настольных приложений на основе Go и VueJS. [Wails](https://github.com/wailsapp/wails) о

###### macOS

Убедитесь, что у вас установлены инструменты командной строки xcode. Для этого выполните следующую команду:
```
xcode-select --
```

###### Debian/Ubuntu

sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev

Дистрибутивы: Debian: 8, 9, 10 | Ubuntu: 16.04, 18.04, 19.04

Также успешно протестировано на: Zorin 15, Parrot 4.7, Linuxmint 19, Elementary 5, Kali, Neon

###### Arch Linux

sudo pacman -S webkit2gtk gtk3

Также успешно протестировано на: Manjaro и ArcoLinux

###### CentOS

sudo yum install webkitgtk3-devel gtk3-devel

CentOS: 6, 7

###### Fedora

sudo yum install webkit2gtk3-devel gtk3-devel

Fedora: 29, 30

###### VoidLinux и VoidLinux-musl

xbps-install gtk+3-devel webkit2gtk-devel

###### Gentoo

sudo emerge gtk+:3 webkit-gtk

###### Windows

Для Windows требуется gcc и связанные инструменты. Рекомендуется загрузить их с [http://tdm-gcc.tdragon.net/download](http://tdm-gcc.tdragon.net/download). После установки вы готовы к работе.

**Убедитесь, что модули Go включены: GO111MODULE=on, и путь к go/bin находится в переменной PATH.**

Установка сводится к выполнению следующей команды:

go get -u github.com/wailsapp/wails/cmd/wails

##### 4. Клонируйте этот репозиторий в ваш GOPATH

git clone https://github.com/grvlle/constellation_wallet.git
```


### Хотите внести свой вклад? Отлично!

Shiwa Wallet использует Wails + Webpack для быстрой разработки интерфейса. Внесите изменение в свой файл и мгновенно увидите обновления!

Откройте ваш терминал и выполните следующие команды.

В директории shiwa_wallet запустите:

$ wails serve

В директории frontend запустите:

$ npm run serve
```

##### Альтернативно:
Если вы хотите скомпилировать кошелек, просто выполните:

wails build

или для отладочной версии:

wails build -d


### Лицензия

Shiwa Wallet распространяется под лицензией MIT.
