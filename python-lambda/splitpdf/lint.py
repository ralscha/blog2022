import subprocess
import sys


def run_lint():
    commands = [
        [sys.executable, "-m", "ruff", "format", "splitpdf", "tests"],
        [sys.executable, "-m", "ruff", "check", "splitpdf", "tests"],
        [sys.executable, "-m", "mypy", "splitpdf", "tests"],
        [sys.executable, "-m", "pytest"],
    ]

    for command in commands:
        result = subprocess.run(command)
        if result.returncode != 0:
            sys.exit(result.returncode)


if __name__ == "__main__":
    run_lint()
