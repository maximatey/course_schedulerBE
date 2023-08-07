# Course Scheduler (Back-end)

Repositori ini mengandung bagian back-end dari Web-App yang berfungsi untuk menjadwalkan mata kuliah yang dapat dipilih untuk mendapatkan nilai maksimal.

Front End: [Click Here](https://github.com/maximatey/course_schedulerFE)
# Tech and Framework

* Go
* MySQL
* Gin

# Penjelasan Dynamic Programming

Dynamic programming adalah pendekatan dalam pemrograman komputer untuk memecahkan masalah yang melibatkan submasalah yang lebih kecil dengan menyimpan solusi submasalah tersebut dan menggabungkannya untuk membangun solusi untuk masalah yang lebih besar. Pendekatan ini memungkinkan kita menghindari perhitungan berulang-ulang yang sama dan secara efisien mengatasi masalah yang kompleks.


Dalam dynamic programming, solusi untuk setiap submasalah dihitung sekali dan disimpan dalam tabel atau struktur data lainnya. Kemudian, solusi submasalah ini digunakan untuk membangun solusi masalah yang lebih besar secara bertahap.

# Analisis Algoritma
Algoritma yang digunakan bersifat rekursif, dimana program akan memeriksa kondisi IP sekumpulan Courses apabila suatu course ditambahkan atau tidak. Algoritma kemudian mengambil sekumpulan courses dengan IP tertinggi.

Proses filtering Courses yang sesuai kriteria terjadi pada method _GetBestValueCourses_, kemudian Courses yang melewati filter tersebut akan dipanggil ke method _findCombinations_.

Pencarian kombinasi Courses dicari secara rekursif dana method _findCombinations_

# Cara Menjalankan Aplikasi

* Clone/Download repositori ini
* Pada root folder, jalankan perintah

```
    docker-compose up -d
```

# Referensi Belajar

* [Strategi Algoritma 2022/2023 - Rinaldi Munir](https://informatika.stei.itb.ac.id/~rinaldi.munir/Stmik/2020-2021/Program-Dinamis-2020-Bagian1.pdf)
* [Dynamic Programming - Geeksforgeeks](https://informatika.stei.itb.ac.id/~rinaldi.munir/Stmik/2020-2021/Program-Dinamis-2020-Bagian1.pdf)

# Catatan

Database yang digunakan program memiliki host `freemysqlhosting.net`, dimana servernya kadang menjadi down sehingga database tidak dapat diakses.
Database juga bersifat temporary, dimana database tidak berlaku selamanya.
