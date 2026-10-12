; факториал: вложенные вызовы, символы, числа и строка
(define (factorial n)
  (if (= n 0)
      1
      (* n (factorial (- n 1)))))

(display "factorial 5")
(factorial 5)
