<!-- Улучшенная совместимость ссылки "Наверх" -->
<a id="readme-top"></a>

  <h1 align="center">WikiLinkExplorer</h1>

  <p align="center">
    Многопоточный поисковик путей в Википедии на Go!
    <br />
    <br />
    <a href="https://github.com/kartmos/WikiLinkExplorer/issues/new?labels=bug&template=bug-report.md">Сообщить об ошибке</a>
    &middot;
    <a href="https://github.com/kartmos/WikiLinkExplorer/issues/new?labels=enhancement&template=feature-request.md">Предложить улучшение</a>
  </p>
</div>

<!-- ОГЛАВЛЕНИЕ -->
<details>
  <summary>Оглавление</summary>
  <ol>
    <li><a href="#о-проекте">О проекте</a></li>
    <li><a href="#возможности">Возможности</a></li>
    <li><a href="#технологии">Технологии</a></li>
    <li><a href="#начало-работы">Начало работы</a></li>
    <li><a href="#использование">Использование</a></li>
    <li><a href="#планы">Планы</a></li>
    <li><a href="#вклад-в-проект">Вклад в проект</a></li>
    <li><a href="#лицензия">Лицензия</a></li>
    <li><a href="#контакты">Контакты</a></li>
  </ol>
</details>

<!-- О ПРОЕКТЕ -->
## О проекте

WikiLinkExplorer - высокопроизводительный веб-краулер для поиска кратчайшего пути между статьями Википедии с использованием конкурентного программирования на Go. Основные особенности:

* Параллельный краулинг с горутинами
* Обработка таймаутов через контексты
* Коммуникация через каналы
* Парсинг HTML с регулярными выражениями

<p align="right">(<a href="#readme-top">наверх</a>)</p>

### Возможности

- Многопоточный поиск
- Настраиваемый таймаут (по умолчанию: 5 минут)
- Визуализация прогресса
- Отслеживание совпадений через уровни статей
- Обработка ошибок с повторами

<!-- НАЧАЛО РАБОТЫ -->
## Начало работы

### Требования

- Go 1.21+
- Интернет-соединение

### Установка

1. Клонировать репозиторий
```sh
   git clone https://github.com/kartmos/WikiLinkExplorer.git
```
<p align="right">(<a href="#readme-top">наверх</a>)</p><!-- ИСПОЛЬЗОВАНИЕ -->

### Использование

1. Переходим в директорию проекта
```sh
cd WikiLinkExplorer
```

2. Собираем бинарный файл
```sh
go build -o wiki-explorer cmd/wikilinkexplorer/main.go
```
3. Вводим параметры запуска
```sh
./wiki-explorer \
  -start="https://en.wikipedia.org/wiki/World" \
  -target="https://en.wikipedia.org/wiki/War" \
  -threads=8 \
  -timeout=10m
```

Пример вывода:
```sh
Matched on level 3:


---> https://en.wikipedia.org/wiki/War
```


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
<!-- БЕЙДЖИ ПРОЕКТА -->
[![Контрибьюторы][contributors-shield]][contributors-url]
[![Форки][forks-shield]][forks-url]
[![Звёзды][stars-shield]][stars-url]
[![Проблемы][issues-shield]][issues-url]

[contributors-shield]: https://img.shields.io/github/contributors/othneildrew/Best-README-Template.svg?style=for-the-badge
[contributors-url]: https://github.com/kartmos/WikiLinkExplorer/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/othneildrew/Best-README-Template.svg?style=for-the-badge
[forks-url]: https://github.com/kartmos/WikiLinkExplorer/network/members
[stars-shield]: https://img.shields.io/github/stars/othneildrew/Best-README-Template.svg?style=for-the-badge
[stars-url]: https://github.com/kartmos/WikiLinkExplorer/stargazers
[issues-shield]: https://img.shields.io/github/issues/othneildrew/Best-README-Template.svg?style=for-the-badge
[issues-url]: https://github.com/kartmos/WikiLinkExplorer/issues
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/kartmos/WikiLinkExplorer/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
[product-screenshot]: images/screenshot.png
[Next.js]: https://img.shields.io/badge/next.js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white
[Next-url]: https://nextjs.org/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Vue.js]: https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D
[Vue-url]: https://vuejs.org/
[Angular.io]: https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white
[Angular-url]: https://angular.io/
[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00
[Svelte-url]: https://svelte.dev/
[Laravel.com]: https://img.shields.io/badge/Laravel-FF2D20?style=for-the-badge&logo=laravel&logoColor=white
[Laravel-url]: https://laravel.com
[Bootstrap.com]: https://img.shields.io/badge/Bootstrap-563D7C?style=for-the-badge&logo=bootstrap&logoColor=white
[Bootstrap-url]: https://getbootstrap.com
[JQuery.com]: https://img.shields.io/badge/jQuery-0769AD?style=for-the-badge&logo=jquery&logoColor=white
[JQuery-url]: https://jquery.com 

