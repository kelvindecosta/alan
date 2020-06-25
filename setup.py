from setuptools import setup, find_packages

with open("README.md", "r") as fh:
    long_description = fh.read()


setup(
    name="alan",
    version="0.4.2",
    description="A programming language for designing Turing machines",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/kelvindecosta/alan",
    author="Kelvin DeCosta",
    author_email="decostakelvin@gmail.com",
    license="MIT",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
    packages=["alan"],
    package_dir={"alan": "src"},
    install_requires=["imageio", "pydot"],
    entry_points={"console_scripts": ["alan = alan.__main__:main",]},
)
