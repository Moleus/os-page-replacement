#import "@local/itmo-report:1.0.0": *

#let variant = 91582

#show: project.with(
  subject: "Операционные системы",
  title: [Лабораторная работа №1. Вариант #variant],
  student: "Соловьев Павел Андреевич",
  group: "P33302",
  teacher: "Макарьев Евгений Юрьевич",
)

#set heading(numbering: none)


#let snip = (name) => raw(read(name), lang: "sql")

#let i = (path, cap) => figure(
  image(path, width: 40%),
  supplement: "",
  numbering: none,
  caption: [ План #cap ],
)


== Задание
#image("./img/task.png")

== Ход работы

Была написана программа на `Go`, рассматривающая данные стратегии, также учитывающая уменьшение и увеличение количества кадров.

#let algorithm_src = read("./pkg/page-replacement/algorithms.go")
#show raw: set text(font: "JuliaMono", size: 0.71em)
#raw(algorithm_src, lang: "Go")

=== Вывод программы
#let fifo_res = read("./results/fifo-7-frames.txt")
#raw(fifo_res, lang: "txt")

#line(length: 100%)

#let lru_res = read("./results/lru-7-frames.txt")
#raw(lru_res, lang: "txt")

#line(length: 100%)

#let opt_res = read("./results/opt-7-frames.txt")
#raw(opt_res, lang: "txt")

==== Количество кадров: 3

#let fifo_res = read("./results/fifo-3-frames.txt")
#raw(fifo_res, lang: "txt")

#line(length: 100%)

#let lru_res = read("./results/lru-3-frames.txt")
#raw(lru_res, lang: "txt")

#line(length: 100%)

#let opt_res = read("./results/opt-3-frames.txt")
#raw(opt_res, lang: "txt")

==== Количество кадров: 14

#let fifo_res = read("./results/fifo-14-frames.txt")
#raw(fifo_res, lang: "txt")

#line(length: 100%)

#let lru_res = read("./results/lru-14-frames.txt")
#raw(lru_res, lang: "txt")

#line(length: 100%)

#let opt_res = read("./results/opt-14-frames.txt")
#raw(opt_res, lang: "txt")

==== 5% страничных сбоев при оптимальном алгоритме

#let fifo_res = read("./results/brute-frames-percents.txt")
#raw(fifo_res, lang: "txt")


=== Ответы на вопросы
- Как изменится количество замен страниц, если увеличить количество кадров в 2 раза?
  уменьшится примерно в #{14 / 6} раза
- А если уменьшить количество кадров в 2 раза?
  увеличится примерно в #{25 / 14} раза
- Сколько должно кадров в памяти, чтобы оптимальный алгоритм давал 5% страничных сбоев?
  19 и более

