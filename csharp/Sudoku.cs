using System;
using System.IO;
using System.Net.Http;
using System.Threading.Tasks;

namespace sudoku
{
    public class Sudoku
    {
        private int[,] values;
        private const string WebSudokuEndpoint = "http://view.websudoku.com/";
        private static HttpClient client = new HttpClient();

        public Sudoku(int[,] values)
        {
            // if the dimensions are not precisely 9x9, assume any value outside
            // the 9x9 grid is null rather than throw an exception on instantiation
            this.values = new int[9, 9];

            if (values.Rank != 2)
                return;

            for (int r = 0; r < Math.Min(9, values.GetLongLength(0)); r++)
                for (int c = 0; c < Math.Min(9, values.GetLongLength(1)); c++)
                    this.values[r, c] = values[r, c];

        }

        public static Sudoku FromTextReader(TextReader reader)
        {
            int[,] values = new int[9, 9];
            int ln = 0;
            while (ln < 9)
            {
                var line = reader.ReadLine();
                if (line == null)
                    break;
                if (line.IndexOf("//") == 0 || line.Length == 0)
                    continue;
                string[] inputs = line.Split(" ", 9);
                if (inputs.Length != 9)
                    throw new Exception($"cannot read row of {inputs.Length} elements");

                for (int i = 0; i < inputs.Length; i++)
                {
                    if (inputs[i] == ".")
                        continue;
                    try
                    {
                        int value = Int32.Parse(inputs[i]);
                        if (value < 1 || value > 9)
                            throw new Exception($"invalid value: cannot read {inputs[i]} in puzzle");
                        values[ln, i] = value;
                    }
                    catch (FormatException)
                    {
                        throw new Exception($"cannot parse value {inputs[i]}");
                    }
                }
                ln++;
            }
            if (ln == 9)
            {
                return new Sudoku(values);
            }
            else
            {
                // all lines valid but puzzle incomplete, so return an exception
                throw new Exception("incomplete puzzle found");
            }
        }

        public static Sudoku FromWebsudoku(int level, int id)
        {
            // spinning up a task and immediately waiting seems fine in this case.
            Task<string> asyncContents = client.GetStringAsync(WebSudokuEndpoint + $"?level={level}&set_id={id}");
            asyncContents.Wait();
            string contents = asyncContents.Result;
            int[,]? puzzle = parsePuzzleHtml(contents);
            if (puzzle == null)
            {
                throw new Exception("error fetching puzzle");
            }
            else
                return new Sudoku(puzzle);
        }

        // Parse out the puzzle from the html. NOTE: this is hacky in that it operates
        // directly on the html string as opposed to parsed html. This project
        // is purely educational. I may remove this in favor of proper parsing later.
        private static int[,]? parsePuzzleHtml(string html)
        {
            int start = html.IndexOf("<TABLE id=\"puzzle_grid\"");
            if (start == -1)
            {
                return null;
            }
            string puzzleTable = html.Substring(start);
            puzzleTable = puzzleTable.Substring(0, puzzleTable.IndexOf("</TABLE>"));

            int[,] puzzle = new int[9, 9];
            int pos = 0;
            for (int row = 0; row < 9; row++)
            {
                pos = puzzleTable.IndexOf("<TR", pos);
                for (int col = 0; col < 9; col++)
                {
                    pos = puzzleTable.IndexOf("<TD", pos);
                    pos = puzzleTable.IndexOf("<INPUT", pos);
                    string cellInput = puzzleTable.Substring(pos, puzzleTable.IndexOf(">", pos) - pos);
                    int valPos = cellInput.IndexOf("VALUE=\"");
                    if (valPos == -1)
                    {
                        puzzle[row, col] = 0;
                    }
                    else
                    {
                        puzzle[row, col] = cellInput.Substring(valPos, 8).ToCharArray()[7] - '0';
                    }
                }
            }

            return puzzle;
        }

        // TODO: consolidate this logic with lambdas to avoid duplication

        // returns true if there are no constraint violation in a row
        private bool isValidRow(int row)
        {
            bool[] seen = new bool[9];
            for (int c = 0; c < 9; c++)
            {
                int value = this.values[row, c];
                if (value != 0 && seen[value - 1])
                    return false;
                if (value != 0)
                    seen[value - 1] = true;
            }
            return true;
        }

        // returns true if there is no constraint violation in a column
        private bool isValidCol(int col)
        {
            bool[] seen = new bool[9];
            for (int r = 0; r < 9; r++)
            {
                int value = this.values[r, col];
                if (value != 0 && seen[value - 1])
                    return false;
                if (value != 0)
                    seen[value - 1] = true;
            }
            return true;
        }

        // returns true if there is no constraint violation in a block
        private bool isValidBlock(int blockRow, int blockCol)
        {
            bool[] seen = new bool[9];

            for (int dr = 0; dr < 3; dr++)
            {
                for (int dc = 0; dc < 3; dc++)
                {
                    int r = blockRow * 3 + dr;
                    int c = blockCol * 3 + dc;

                    int value = this.values[r, c];
                    if (value != 0 && seen[value - 1])
                        return false;
                    if (value != 0)
                        seen[value - 1] = true;
                }
            }
            return true;
        }

        public void Print(TextWriter writer)
        {
            for (int row = 0; row < 9; row++)
            {
                for (int col = 0; col < 9; col++)
                {
                    if (this.values[row, col] == 0)
                        writer.Write(".");
                    else
                        writer.Write((char)('0' + this.values[row, col]));
                    if (col != 8)
                        writer.Write(" ");
                }
                writer.WriteLine();
            }
        }

        public bool Solve()
        {
            return solve(0);
        }

        // brute force recursive search for solution
        private bool solve(int iter)
        {
            // checking the 82nd iteration in a sudoku puzzle of 81 boxes implies success
            if (iter == 81)
                return true;

            int row = iter / 9;
            int col = iter % 9;

            // if there's already an assigned number, skip by recursing
            if (this.values[row, col] != 0)
            {
                return solve(iter + 1);
            }
            else
            {
                // iterate through all the possible numbers and recurse on non-violations
                for (int n = 1; n <= 9; n++)
                {
                    this.values[row, col] = n;
                    if (isValidCol(col) &&
                     isValidRow(row) &&
                     isValidBlock(row / 3, col / 3) &&
                     solve(iter + 1))
                    {
                        return true;
                    }
                    else
                    {
                        this.values[row, col] = 0; // reset so that failures don't modify the original input array
                    }
                }
            }
            return false;
        }

    }
}