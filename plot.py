import plotly.graph_objects as go
import pandas as pd


def pick(df: pd.DataFrame) -> pd.DataFrame:
    low = 10
    high = df['value'].quantile(0.95)
    return df.query('@low < value < @high')


def convert(df: pd.DataFrame) -> (pd.DataFrame, list):
    _, labels = df['source'].factorize()
    idx = range(len(labels))
    dic = dict(zip(labels, idx))
    return df.replace(dic), labels.to_list()


def main():
    df = pd.read_csv('data.csv')
    # df = pick(df)
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
