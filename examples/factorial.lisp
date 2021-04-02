;; this implementation is tail-optimised
(def! factorial
  (fn* [k]
    (let* [factorial-inner
	  (fn* [x acc]
		(let* [decremented (- x 1) nacc (* x acc)]
		  (cond
			(< x 0) (throw "negative values are not allowed")
			(= x 0) 1
			(= decremented 1) nacc
			"else" (factorial-inner decremented nacc))))]
      (factorial-inner k 1))))

(def! ! factorial)

(println "The factorial of -1 is:" (! -1))
(println "The factorial of 0 is:" (! 0))
(println "The factorial of 1 is:" (! 1))
(println "The factorial of 2 is:" (! 2))
(println "The factorial of 3 is:" (! 3))
(println "The factorial of 4 is:" (! 4))
(println "The factorial of 5 is:" (! 5))
(println "The factorial of 1000 is:" (! 1000))
