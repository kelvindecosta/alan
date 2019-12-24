import imageio
import re
from pydot import Dot, Edge, Node, Subgraph

class Machine:
    """
    The abstract class for a Turing Machine.
    """

    def __init__(self):
        """
        Initializes the Machine object.
        """
        self._symbols = set()
        self._blank_symbol = None
        self._states = set()
        self._start_state = None
        self._end_states = set()
        self._transitions = {}

        self._current_state = None
        self._tape = None
        self._head = None


    def _set_symbol(self, symbol, blank=False):
        """
        Adds a symbol to the set of symbols of the Machine.
        Sets the blank symbol only if it is not already set.

        Arguments:
            symbol {str} -- symbol

        Keyword Arguments:
            blank {bool} -- whether the symbol must be set as the blank symbol (default: {False})
        """
        self._symbols.add(symbol)

        try:
            assert self._blank_symbol == None or not blank
            if blank:
                self._blank_symbol = symbol
        except:
            raise Exception(f"Machine got blank symbol '{symbol}' which is already set to '{self._blank_symbol}'")


    def _set_state(self, state, start=False, end=False):
        """
        Adds a state to the set of states of the Machine.
        Sets the start state only if it is not already set.

        Arguments:
            state {str} -- state

        Keyword Arguments:
            start {bool} -- whether the symbol must be set as the start state (default: {False})
            end {bool} -- whether the symbol must be set as the end state (default: {False})
        """
        self._states.add(state)
        if end:
            self._end_states.add(state)

        try:
            assert self._start_state == None or not start
            if start:
                self._start_state = state
        except:
            raise Exception(f"Machine got start state '{state}' which is already set to '{self._start_state}'")


    def _set_transition(self, current_state, current_symbol, next_symbol, direction, next_state):
        """
        Adds a transition to the set of transitions of the Machine.

        Arguments:
            current_state {str} -- current State
            current_symbol {str} -- current Symbol
            next_symbol {str} -- next Symbol
            direction {bool} -- direction
            next_state {str} -- next State
        """
        self._set_symbol(current_symbol)
        self._set_symbol(next_symbol)
        self._set_state(current_state)
        self._set_state(next_state)

        if self._transitions.get(current_state) is None:
            self._transitions[current_state] = {}

        self._transitions[current_state][current_symbol] = (next_symbol, direction, next_state)


    def parse(self, definition):
        """
        Parses a definition for a Machine object.

        Arguments:
            definition {str} -- definition
        """
        comment_re = re.compile(r"(#.*)")
        state_re = re.compile(r"^([a-zA-Z_][a-zA-Z0-9_]*)([.*])?$")
        transition_re = re.compile(r"^'(.)'\s+'(.)'\s+([<>])\s+([a-zA-Z_][a-zA-Z0-9_]*)$")

        lines = list(filter(lambda x : len(x) > 0, map(lambda x: comment_re.sub("", x).strip(), definition.split("\n"))))
        self._set_symbol(lines.pop(0)[1], True)

        scope_state = None
        for l in lines:
            try:
                g = state_re.match(l).groups()
                scope_state = g[0]
                special = g[1] if g[1] is not None else ""
                self._set_state(scope_state, "*" in special, "." in special)
            except:
                pass

            try:
                g = transition_re.match(l).groups()
                self._set_transition(scope_state, g[0], g[1], g[2] is ">", g[3])
            except:
                pass


    def reset(self, tape):
        """
        Resets the tape, head and state of the Machine.

        Arguments:
            tape {[type]} -- string of tape symbols
        """
        self._current_state = self._start_state
        self._tape = list(tape)
        self._head = 0


    def step(self):
        """
        Performs one Machine step.

        Returns:
            {bool} -- whether or not the Machine halted
        """
        try:
            current_symbol = self._tape[self._head]
            next_symbol, direction, self._current_state = self._transitions.get(self._current_state).get(current_symbol)
        except:
            return True

        self._tape[self._head] = next_symbol
        self._head += 1 if direction else -1

        if self._head < 0:
            self._tape.insert(0, self._blank_symbol)
            self._head = 0
        elif self._head >= len(self._tape):
            self._tape.append(self._blank_symbol)
            self._head = len(self._tape) - 1

        return False


    def run(self, tape, max_steps=200, animate=False, **kwargs):
        """
        Performs a computation by the Machine on a tape

        Arguments:
            tape {str} -- string of tape symbols

        Keyword Arguments:
            max_steps {int} -- maximum steps before forcing halt (default: {200})
            animate {bool} -- whether to call animation functionality (default: {False})

        Returns:
            ({bool}, {bool}, {str}) -- whether it halted, whether it accepted, final tape
        """
        self.reset(tape)
        halt = False

        if animate:
            try:
                assert kwargs.get("filename") is not None
            except:
                raise Exception("Specify a filename to save the animation")
            images = []


        for _ in range(max_steps):
            if animate:
                images.append(imageio.imread(self.graph(False).create(prog="dot", format="png")))
                images.append(imageio.imread(self.graph(True).create(prog="dot", format="png")))

            halt = self.step()
            if halt:
                break

        if animate:
            imageio.mimsave(kwargs.get("filename"), images, fps=kwargs.get("fps"))

        return halt, halt and self._current_state in self._end_states, "".join(self._tape).strip(self._blank_symbol)


    def graph(self, context=None, **kwargs):
        """
        Returns a graph object

        Keyword Arguments:
            context {bool|None} -- controls the animations (default: {None})

        Returns:
            [pydot.Dot] -- graph object
        """
        graph = Dot(graph_type="digraph", rankdir=("LR" if context is None else "TB"))
        machine_graph =  Subgraph(graph_name="cluster_machine", graph_type="digraph", label="MACHINE")

        for current_state in sorted(self._states):
            node_args = {}
            shape = "circle"
            shape = ("double" if current_state in self._end_states else "") + shape
            node_args["shape"] = shape

            if context is not None and current_state == self._current_state:
                node_args["fillcolor"] = "cyan"
                node_args["style"] = "filled"

            machine_graph.add_node(Node(current_state, **node_args))

        machine_graph.add_node(Node("0", shape="point"))
        machine_graph.add_edge(Edge("0", self._start_state))

        for current_state in self._states:
            transitions = self._transitions.get(current_state)
            if transitions:
                for current_symbol in transitions:
                    next_symbol, direction, next_state = transitions.get(current_symbol)
                    label = f"'{current_symbol}' '{next_symbol}' {'R' if direction else 'L'}"

                    edge_args = {}
                    if context and current_state == self._current_state and current_symbol == self._tape[self._head]:
                        edge_args["color"] = "cyan"
                    machine_graph.add_edge(Edge(current_state, next_state, label=label, **edge_args))

        graph.add_subgraph(machine_graph)
        if context is not None:
            tape_graph = Subgraph(graph_name="cluster_tape", graph_type="digraph", label="TAPE")
            tape = []
            for index in range(-4 + self._head, 5 + self._head):
                tape.append(f"<t{index}> {self._tape[index] if 0 <= index < len(self._tape) else self._blank_symbol}")

            tape_graph.add_node(Node("tape", label="|".join(tape), shape="record"))
            tape_graph.add_node(Node("t0", shape="point"))
            tape_graph.add_edge(Edge("t0", f"tape:t{self._head}"))
            graph.add_subgraph(tape_graph)

        if kwargs.get("filename"):
            graph.write(kwargs.get("filename"), format="png")
            return f"Graph saved to {kwargs.get('filename')}"

        return graph
