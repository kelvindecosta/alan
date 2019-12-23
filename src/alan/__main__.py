from .interface import parse_args
from .machine import Machine


def main():
    args = vars(parse_args())
    command = args.pop("command")
    definition = args.pop("definition")

    m = Machine()
    with open(definition, "r") as f:
        m.parse(f.read())

    result = getattr(m, command)(**args)

    if command == "run":
        if all(result[:2]):
            print("Accepted")
        elif result[0]:
            print("Rejected")
        else:
            print("Undecidable")
        print(f"Initial Tape : {args.get('tape')}\nFinal Tape   : {result[2]}")
    else:
        print(result)

if __name__ == "__main__":
    main()
