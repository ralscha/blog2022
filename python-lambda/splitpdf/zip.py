import subprocess
import sys
import os
import shutil
import zipfile


def run_command(command: str) -> None:
    result = subprocess.run(command, shell=True)
    if result.returncode != 0:
        sys.exit(result.returncode)


def zip_directory(directory_path: str, zip_path: str) -> None:
    with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zipf:
        for root, dirs, files in os.walk(directory_path):
            for file in files:
                file_path = os.path.join(root, file)
                arcname = os.path.relpath(file_path, directory_path)
                zipf.write(file_path, arcname)


def copy_directory(src: str, dst: str) -> None:
    for item in os.listdir(src):
        s = os.path.join(src, item)
        d = os.path.join(dst, item)
        if os.path.isdir(s):
            if item == "__pycache__":
                continue
            shutil.copytree(s, d, ignore=shutil.ignore_patterns("__pycache__"))
        else:
            shutil.copy2(s, d)


def run_zip():
    run_command("poetry install --only main --sync")
    shutil.rmtree("dist", ignore_errors=True)
    os.makedirs("dist/lambda-package/splitpdf", exist_ok=True)

    copy_directory(".venv/lib/site-packages", "dist/lambda-package")
    copy_directory("splitpdf", "dist/lambda-package/splitpdf")

    zip_directory("dist/lambda-package", "dist/lambda.zip")
