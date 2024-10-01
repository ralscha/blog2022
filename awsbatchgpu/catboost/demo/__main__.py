import os

import boto3
import pandas as pd
from catboost import CatBoostClassifier, Pool
from mypy_boto3_s3.client import S3Client
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder


def main():
  url = "https://archive.ics.uci.edu/ml/machine-learning-databases/adult/adult.data"
  column_names = [
    "age",
    "workclass",
    "fnlwgt",
    "education",
    "education-num",
    "marital-status",
    "occupation",
    "relationship",
    "race",
    "sex",
    "capital-gain",
    "capital-loss",
    "hours-per-week",
    "native-country",
    "income",
  ]
  data = pd.read_csv(url, names=column_names, skipinitialspace=True)

  le = LabelEncoder()
  for column in data.select_dtypes(include=["object"]):
    data[column] = le.fit_transform(data[column])

  X = data.drop("income", axis=1)
  y = data["income"]

  x_train, x_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=2306
  )

  cat_features = [
    "workclass",
    "education",
    "marital-status",
    "occupation",
    "relationship",
    "race",
    "sex",
    "native-country",
  ]

  train_pool = Pool(x_train, y_train, cat_features=cat_features)
  test_pool = Pool(x_test, y_test, cat_features=cat_features)
  model = CatBoostClassifier(
    iterations=5000,
    task_type="GPU",
    depth=8,
    loss_function="Logloss",
    verbose=False,
  )

  model.fit(train_pool, eval_set=test_pool, verbose=500)
  model.save_model("model.cbm")

  s3: S3Client = boto3.client("s3")
  s3.upload_file("model.cbm", os.environ["OUTPUT_BUCKET"], "model.cbm")


if __name__ == "__main__":
  main()
