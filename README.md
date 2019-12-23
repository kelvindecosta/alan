# Alan

<p align=center>

  <img src="https://raw.githubusercontent.com/kelvindecosta/alan/master/assets/readme/logo.png" height="200px"/>

  <br>
  <span>A programming language for <em>designing</em> Turing machines.</span>
  <br>
  <a target="_blank" href="https://www.python.org/downloads/" title="Python version"><img src="https://img.shields.io/badge/python-%3E=_3.6-green.svg"></a>
  <a target="_blank" href="LICENSE" title="License: MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
  <a target="_blank" href="https://pypi.python.org/pypi/alan/"><img alt="pypi package" src="https://badge.fury.io/py/alan.svg"></a>
</p>

<p align="center">
  <a href="#walkthrough">Walkthrough</a>
  &nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#installation">Installation</a>
  &nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="https://github.com/kelvindecosta/alan/wiki">Wiki</a>
  &nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#citation">Citation</a>

</p>

## Installation

```bash
pip install alan
```

## Walkthrough

This section describes a workflow.

For an in-depth guide navigate to the [Wiki](https://github.com/kelvindecosta/alan/wiki).
Here are some useful links:

*   [Definition](https://github.com/kelvindecosta/alan/wiki/Definition)
*   [Syntax](https://github.com/kelvindecosta/alan/wiki/Syntax)
*   [Interface](https://github.com/kelvindecosta/alan/wiki/Interface)

Consider the following example, the definition for a Turing machine that accepts all binary strings that are palindromic:

```
# This is a definition of a Turing Machine that accepts binary strings that are palindromes
' '
A*
    'X' 'X' < A
    'Y' 'Y' < A
    '0' 'X' > B
    '1' 'Y' > F
    ' ' ' ' > G
B                   # Starting with 0
    '0' '0' > B
    '1' '1' > B
    ' ' ' ' < C
    'X' 'X' < C
    'Y' 'Y' < C
F                   # Starting with 1
    '0' '0' > F
    '1' '1' > F
    ' ' ' ' < E
    'X' 'X' < E
    'Y' 'Y' < E
C
    '0' 'X' < D
    'X' 'X' < D
E
    '1' 'Y' < D
    'Y' 'Y' < D
D
    '0' '0' < D
    '1' '1' < D
    ' ' ' ' > A
    'X' 'X' > A
    'Y' 'Y' > A
G.
    'X' '0' > G
    'Y' '1' > G
```

Graph the machine:

```bash
alan graph examples/binary-palindrome.aln -f assets/readme/binary-palindrome.png
```

<p align="center"><img src="https://raw.githubusercontent.com/kelvindecosta/alan/master/assets/readme/binary-palindrome.png"></p>

Run the machine on some inputs:

*   ```bash
    alan run examples/binary-palindrome.aln 101
    ```

    ```
    Accepted
    Initial Tape : 101
    Final Tape   : 10
    ```

*   ```bash
    alan run examples/binary-palindrome.aln 1010
    ```

    ```
    Rejected
    Initial Tape : 1010
    Final Tape   : Y010
    ```

Animate the computation on some inputs:

*   ```bash
    alan run examples/binary-palindrome.aln 101 -a -f assets/readme/binary-palindrome-accepted.gif
    ```

    ![Animation of accepted input](https://raw.githubusercontent.com/kelvindecosta/alan/master/assets/readme/binary-palindrome-accepted.gif)

*   ```bash
    alan run examples/binary-palindrome.aln 1010 -a -f assets/readme/binary-palindrome-rejected.gif
    ```

    ![Animation of rejected input](https://raw.githubusercontent.com/kelvindecosta/alan/master/assets/readme/binary-palindrome-rejected.gif)

## Citation

If you use this implementation in your work, please cite the following:

```
@misc{decosta2019alan,
    author = {Kelvin DeCosta},
    title = {Alan},
    year = {2019},
    howpublished = {\url{https://github.com/kelvindecosta/alan}},
}
```
