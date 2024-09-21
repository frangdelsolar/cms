import os


class Colors:
    GREEN = "\033[0;32m"
    RED = "\033[0;31m"
    BLUE = "\033[0;34m"
    CYAN = "\033[0;36m"
    YELLOW = "\033[0;33m"
    WHITE = "\033[0m"
    BLACK = "\033[0;30m"
    BOLD = "\033[1m"
    ITALIC = "\033[3m"

log_colors = {
    "FIXME": Colors.RED,
    "TODO": Colors.YELLOW
}


results_count = {
    "FIXME": 0,
    "TODO": 0,
}


def read_go_files(folder_path: str) -> list[str]:
    """
    Reads all .go files in the given folder and its subfolders and returns a list of the paths to the files.
    """
    files = []
    for root, dirs, filenames in os.walk(folder_path):
        for filename in filenames:
            if filename.endswith(".go"):
                files.append(os.path.join(root, filename))
    return files


def find_comments(file_path: str) -> None:
    """
    Reads all lines from the given file and searches for FIXME and TODO comments.

    :param file_path: The path to the file to read from.
    """
    print(f"{Colors.BLACK}Checking: {file_path}")

    with open(file_path, 'r') as file:
        lines = file.readlines()

    for line_number, line in enumerate(lines, 1):  # Use enumerate for line numbers
        if line.strip().startswith("//"):  # Check for line comments
            comment = line.strip()[2:]  # Extract comment text after "//"
            for comment_type in ("FIXME", "TODO"):
                if comment_type.lower() in comment.lower():  # Case-insensitive check
                    results_count[comment_type] += 1
                    print(f"{log_colors[comment_type]}{file_path}:{line_number} -> {comment.strip()}")


def present_results():
    """
    Prints the results of the comment search to the console.
    """

    print("\n\nResults:")
    for comment_type, count in results_count.items():
        color = Colors.YELLOW if comment_type == "FIXME" else Colors.GREEN
        print(f"{color}{comment_type}: {count}")

    if results_count["FIXME"] > 0:
        print(f"{Colors.RED}Error: Found FIXME comments. Please address them before continuing.")
        exit(1)  # Exit with non-zero code to indicate an error


def main():
    folder_path = "../builder"
    files = read_go_files(folder_path)

    for file in files:
        find_comments(file)

    present_results()


if __name__ == "__main__":
    main()