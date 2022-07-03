<div id="top"></div>

# Ya Online Judge

<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache License][license-shield]][license-url]
![Go version (master)][gover-shield]
![GitHub tag (latest by date)][tag-shield]
![GitHub code size in bytes][codesize-shield]



<div align="center">
<!--
  <a href="https://github.com/sshwy/yaoj-core">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>
-->
</div>

Core packages for YaOJ.

<!-- TABLE OF CONTENTS -->
<details open>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <!-- <li><a href="#acknowledgments">Acknowledgments</a></li> -->
  </ol>
</details>

## About The Project

YaOJ, namely Yet Another Online Judge, is coming as a brand new OJ with an excellent framework for easily configured custom testing. This repo contains all core packages of YaOJ, which is for customized problem testing.

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

### Prerequisites

The go toolkit command.

### Installation

1. Clone the repo

   ```sh
   git clone --recursive https://github.com/sshwy/yaoj-core
   ```

2. Generate necessary files

   ```sh
   cd yaoj-core
   go generate
   ```
3. Build Apps

  ```sh
  go build ./cmd/migrator
  go build ./cmd/judgeserver
  ```

4. Happy developing!

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

<a href="https://pkg.go.dev/github.com/sshwy/yaoj-core@master"><img src="https://pkg.go.dev/badge/github.com/sshwy/yaoj-core.svg" alt="Go Reference"></a>

_For more examples, please refer to `*_test.go` files_

<p align="right">(<a href="#top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'feat(scope): add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>

## License

Distributed under the Apache-2.0 License. See `LICENSE` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

## Contact

Project Link: [https://github.com/sshwy/yaoj-core](https://github.com/sshwy/yaoj-core)

<p align="right">(<a href="#top">back to top</a>)</p>

<!--
## Acknowledgments

* []()
* []()
* []()

<p align="right">(<a href="#top">back to top</a>)</p>
-->

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/sshwy/yaoj-core.svg
[contributors-url]: https://github.com/sshwy/yaoj-core/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/sshwy/yaoj-core.svg
[forks-url]: https://github.com/sshwy/yaoj-core/network/members
[stars-shield]: https://img.shields.io/github/stars/sshwy/yaoj-core.svg
[stars-url]: https://github.com/sshwy/yaoj-core/stargazers
[issues-shield]: https://img.shields.io/github/issues/sshwy/yaoj-core.svg
[issues-url]: https://github.com/sshwy/yaoj-core/issues
[license-shield]: https://img.shields.io/github/license/sshwy/yaoj-core.svg
[license-url]: https://github.com/sshwy/yaoj-core/blob/master/LICENSE
[gover-shield]: https://img.shields.io/github/go-mod/go-version/sshwy/yaoj-core/master?filename=go.mod
[tag-shield]: https://img.shields.io/github/v/tag/sshwy/yaoj-core?label=latest%20tag
[codesize-shield]: https://img.shields.io/github/languages/code-size/sshwy/yaoj-core