Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

Ответ:
```
<nil>
false
```

Несмотря на то, что базовая структура является равна `nil`, сам интерфейс не является `nil`. Интерфейс равен `nil`, только если и тип, и значение равны `nil`.

## Устройство интерфейсов
https://github.com/teh-cmc/go-internals/blob/master/chapter2_interfaces/README.md

`iface` — это корневой тип, представляющий интерфейс рантайма ([src/runtime/runtime2.go](https://github.com/golang/go/blob/e822b1e26e20ef1c76672c0b77b0fd8a97a1fe84/src/runtime/runtime2.go#L202)).
Его определение выглядит так:

```go
type iface struct { // 16 bytes on a 64bit arch
    tab  *itab
    data unsafe.Pointer
}
```

Таким образом, интерфейс представляет собой очень простую структуру, которая содержит 2 указателя:
- `tab` содержит адрес объекта `itab`, который описывает как тип интерфейса, так и тип данных, на которые он указывает.
- `data` — это `unsafe` указатель на значение, хранящееся в интерфейсе.

Поскольку интерфейсы могут содержать только указатели, любое конкретное значение, которое мы переносим в интерфейс, должно иметь адрес.

**Структура `itab`**

`itab` определяется следующим образом ([src/runtime/runtime2.go](https://github.com/golang/go/blob/e822b1e26e20ef1c76672c0b77b0fd8a97a1fe84/src/runtime/runtime2.go#L902=)):

```go
type itab struct { // 40 bytes on a 64bit arch
    inter *interfacetype
    _type *_type
    hash  uint32 // copy of _type.hash. Used for type switches.
    _     [4]byte
    fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
```

`Itab` — это сердце и мозг интерфейса.

Во-первых, он имеет поле `_type`, которое является внутренним представлением любого типа Go в рантайме.
`_type` описывает все аспекты типа: его имя, его характеристики (например, размер, выравнивание...) и, в некоторой степени, даже его поведение (например, сравнение, хеширование...)!
В этом случае поле `_type` описывает тип значения, хранящегося в интерфейсе, то есть значение, на которое указывает указатель данных.

Во-вторых, мы находим указатель на `interfacetype`, который представляет собой просто оболочку вокруг _type с некоторой дополнительной информацией, специфичной для интерфейсов.
Как и следовало ожидать, поле `inter` описывает тип самого интерфейса.

Наконец, массив `fun` содержит указатели на функции, составляющие виртуальную/диспетчерскую таблицу интерфейса.
Обратите внимание на комментарий, в котором говорится `// variable sized`, что означает, что размер, с которым объявлен этот массив, не имеет значения.
Компилятор отвечает за выделение памяти, поддерживающей этот массив, и делает это независимо от указанного здесь размера. Точно так же среда выполнения всегда обращается к этому массиву с помощью `unsafe` указателей, поэтому проверка границ здесь не применяется.

`Itab` вычисляется во время выполнения. Рантайм вычисляет `Itab`, ища каждый метод, указанный в таблице методов типа интерфейса, в таблице методов конкретного типа. Рантайм кэширует itable после его создания, так что это соответствие нужно вычислить только один раз.

**Структура `_type`**

Структура `_type` дает полное описание типа Go. Она определяется следующим образом ([src/runtime/type.go](https://github.com/golang/go/blob/e822b1e26e20ef1c76672c0b77b0fd8a97a1fe84/src/runtime/type.go#L35)):

```go
type _type struct { // 48 bytes on a 64bit arch
    size       uintptr
    ptrdata    uintptr // size of memory prefix holding all pointers
    hash       uint32
    tflag      tflag
    align      uint8
    fieldalign uint8
    kind       uint8
    alg        *typeAlg
    // gcdata stores the GC type data for the garbage collector.
    // If the KindGCProg bit is set in kind, gcdata is a GC program.
    // Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
    gcdata    *byte
    str       nameOff
    ptrToThis typeOff
}
```

Типы `nameOff` и `typeOff` представляют собой смещения `int32` в метаданные, встроенные компоновщиком в окончательный исполняемый файл. Эти метаданные загружаются в структуры данных `runtime.moduledata` во время выполнения ([src/runtime/symtab.go](https://github.com/golang/go/blob/e822b1e26e20ef1c76672c0b77b0fd8a97a1fe84/src/runtime/type.go#L35)).

**Структура `interfacetype`**

Структура `interfacetype` ([src/runtime/type.go](https://github.com/golang/go/blob/e822b1e26e20ef1c76672c0b77b0fd8a97a1fe84/src/runtime/type.go#L350)):

```go
type interfacetype struct { // 80 bytes on a 64bit arch
    typ     _type
    pkgpath name
    mhdr    []imethod
}

type imethod struct {
    name nameOff
    ityp typeOff
}
```

`interfacetype` — это просто оболочка вокруг `_type` с некоторыми дополнительными метаданными, специфичными для интерфейса.

В текущей реализации эти метаданные в основном состоят из списка смещений, указывающих на соответствующие имена и типы методов, предоставляемых интерфейсом (`[]imethod`).

Вот краткий обзор того, как выглядит интерфейс `iface`, когда он представлен со всеми встроенными подтипами:

```go
type iface struct { // `iface`
    tab *struct { // `itab`
        inter *struct { // `interfacetype`
            typ struct { // `_type`
                size       uintptr
                ptrdata    uintptr
                hash       uint32
                tflag      tflag
                align      uint8
                fieldalign uint8
                kind       uint8
                alg        *typeAlg
                gcdata     *byte
                str        nameOff
                ptrToThis  typeOff
            }
            pkgpath name
            mhdr    []struct { // `imethod`
                name nameOff
                ityp typeOff
            }
        }
        _type *struct { // `_type`
            size       uintptr
            ptrdata    uintptr
            hash       uint32
            tflag      tflag
            align      uint8
            fieldalign uint8
            kind       uint8
            alg        *typeAlg
            gcdata     *byte
            str        nameOff
            ptrToThis  typeOff
        }
        hash uint32
        _    [4]byte
        fun  [1]uintptr
    }
    data unsafe.Pointer
}
```

**Пустой интерфейс**

Структура данных для пустого интерфейса — это `iface` без `itab`. Тому есть две причины:

- Поскольку в пустом интерфейсе нет методов, все, что связано с динамической диспетчеризацией, можно смело выкинуть из структуры данных.
- Когда виртуальная таблица исчезла, тип самого пустого интерфейса, не путать с типом данных, которые он содержит, всегда один и тот же.

`eface` — это корневой тип, представляющий пустой интерфейс в рантайме ([src/runtime/runtime2.go](https://github.com/golang/go/blob/e822b1e26e20ef1c76672c0b77b0fd8a97a1fe84/src/runtime/runtime2.go#L207)).
Его определение выглядит так:

```go
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
```

Где `_type` содержит информацию о типе значения, на которое указывают данные.
`itab` был полностью удален.

В то время как пустой интерфейс может просто повторно использовать структуру данных `iface` (в конце концов, это надмножество `eface`), среда выполнения предпочитает различать их по двум основным причинам: эффективность использования пространства и ясность кода.

Ссылки:

- https://github.com/teh-cmc/go-internals/blob/master/chapter2_interfaces/README.md
- https://research.swtch.com/interfaces
- https://habr.com/ru/post/276981/