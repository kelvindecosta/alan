from argparse import ArgumentParser


def parse_args():
    """Parses arguments
    """
    parser = ArgumentParser(description="A programming langauge for designing Turing Machines", prog="alan")
    commands = parser.add_subparsers(dest="command")

    # p_run
    p_run = commands.add_parser("run", help="run machine on tape")
    p_run.add_argument("definition", type=str, help="path to definition file")
    p_run.add_argument("tape", type=str, help="string of tape symbols")
    p_run.add_argument("-s", "--max-steps", type=int, help="maximum steps before forcing halt", default=200)
    p_run.add_argument("-a", "--animate", action="store_true")
    p_run.add_argument("-f", "--filename", type=str, help="path to save animation")
    p_run.add_argument("-r", "--fps", type=int, help="animation fps", default=1)

    # p_graph
    p_graph = commands.add_parser("graph", help="graph machine")
    p_graph.add_argument("definition", type=str, help="path to definition file")
    p_graph.add_argument("-f", "--filename", type=str, help="path to save graph")

    args = parser.parse_args()

    if args.command is None:
        parser.print_help()
        exit()

    return args
