import sys


def is_number(s):
    try:
        int(s, 0)
        return True
    except ValueError:
        return False


def main():
    if len(sys.argv) < 2:
        return

    try:
        with open(sys.argv[1], "r") as f:
            lines = f.readlines()
    except Exception as e:
        sys.stderr.write(f"Error reading file: {e}\n")
        sys.exit(1)

    in_forth = False
    for line in lines:
        raw_line = line.rstrip("\n")
        trimmed = raw_line.strip()

        # Check for start of Forth definition: ": name" at beginning of trimmed line
        if not in_forth and trimmed.startswith(": "):
            in_forth = True
            parts = trimmed.split()
            if len(parts) < 2:
                # Should not happen with well-formed Forth, but preserve
                print(raw_line)
                in_forth = False
                continue

            name = parts[1]
            print(f'defcode "{name}", {len(name)}')
            print("    ENTER")

            # Process any tokens following the name on the same line
            for i in range(2, len(parts)):
                t = parts[i]
                if t.startswith("@"):
                    break  # Assembler comment
                if t == ";":
                    print("    EXIT")
                    in_forth = False
                    break
                if is_number(t):
                    print(f"    LIT {t}")
                else:
                    print(f'bl "{t}"')
            continue

        if in_forth:
            if not trimmed:
                print()
                continue

            parts = trimmed.split()
            for t in parts:
                if t.startswith("@"):
                    # Ignore the rest of the line as it's an assembly comment
                    break
                if t == ";":
                    print("    EXIT")
                    in_forth = False
                    break
                if is_number(t):
                    print(f"    LIT {t}")
                else:
                    print(f'bl "{t}"')
            continue

        # Outside of Forth definitions, preserve original assembly lines exactly
        print(raw_line)


if __name__ == "__main__":
    main()
