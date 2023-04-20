import plotly.graph_objects as go
import pandas as pd
import sys

HEADER_SRC = 'source'
HEADER_DST = 'target'
HEADER_CNT = 'value'


def pick(df: pd.DataFrame) -> pd.DataFrame:
    low = df[HEADER_CNT].quantile(0.25)
    high = df[HEADER_CNT].quantile(1.00)
    return df.query(f'@low < {HEADER_CNT} < @high')


def convert(df: pd.DataFrame) -> (pd.DataFrame, list):
    _, labels = df[HEADER_SRC].factorize()
    idx = range(len(labels))
    dic = dict(zip(labels, idx))
    return df.replace(dic), labels.to_list()


def main():
    df = pd.read_csv(sys.stdin)
    df = pick(df)
    df, labels = convert(df)

    fig = go.Figure(data=[go.Sankey(
        node=dict(
            pad=15,
            thickness=20,
            line=dict(color="black", width=0.5),
            label=labels,
            color="blue"
        ),
        link=df.to_dict(orient='list'))])

    fig.show()


main()
