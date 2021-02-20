using System;
using System.Collections.Generic;

namespace sudoku
{
    class Program
    {
        const string helpString = @"Usage of sudoku:
  -dl
        download sudokus
  -id int
        solve puzzle with id 'id' (default -1)
  -level int
        solve puzzle of difficulty 'level' (1, 2, 3, or 4)
  -n int
        download 'n' number of random puzzles";

        static bool getBoolFlag(string[] args, string flag)
        {
            for (int a = 0; a < args.Length; a++)
            {
                if (args[a] == $"-{flag}")
                    return true;
            }
            return false;
        }
        static int? getIntFlag(string[] args, string flag)
        {
            int? val = null;
            for (int a = 0; a < args.Length - 1; a++)
            {
                if (args[a] == $"-{flag}")
                {
                    if (a + 1 < args.Length)
                    {
                        string input = args[a + 1];
                        try
                        {
                            val = Int32.Parse(input);
                        }
                        catch (FormatException)
                        {
                            val = null;
                        }
                    }
                    else
                    {
                        break;
                    }
                }
            }
            return val;
        }

        static (int, int) parseLevelAndId(string[] args)
        {
            int? level = getIntFlag(args, "level");
            int? id = getIntFlag(args, "id");

            if (level == null)
            {
                Console.Error.WriteLine("missing required 'level' param");
            }

            if (id == null)
            {
                Console.Error.WriteLine("missing required 'id' param");
            }

            if (id == null || level == null)
                Environment.Exit(1);

            return (level ?? 0, id ?? 0);
        }

        static void runDownload(string[] args)
        {
            int level, id;
            if (getIntFlag(args, "level") != null || getIntFlag(args, "id") != null)
            {
                (level, id) = parseLevelAndId(args);
                (Sudoku.FromWebsudoku(level, id)).Print(Console.Out);
            }
            else
            {
                int n = getIntFlag(args, "n") ?? 0;
                var rand = new Random();
                for (int i = 0; i < n; i++)
                {
                    level = rand.Next(4) + 1;
                    id = rand.Next(short.MaxValue * 2);
                    Sudoku puzzle = Sudoku.FromWebsudoku(level, id);
                    Console.Out.WriteLine("//fetched from http://www.websudoku.com");
                    Console.Out.WriteLine($"//level: {level}, id: {id}");
                    puzzle.Print(Console.Out);
                    if (i < n - 1)
                    {
                        Console.Out.WriteLine();
                    }
                }
            }

        }

        static void runSolve(string[] args)
        {
            List<Sudoku> puzzles = new List<Sudoku>();
            // require level and id if and only if one already exists
            if (getIntFlag(args, "level") != null || getIntFlag(args, "id") != null)
            {
                int level, id;
                (level, id) = parseLevelAndId(args);
                puzzles.Add(Sudoku.FromWebsudoku(level, id));
            }
            else
            { // otherwise read puzzles from stdin
                while (Console.In.Peek() != -1)
                    puzzles.Add(Sudoku.FromTextReader(Console.In));
            }

            for (int p = 0; p < puzzles.Count; p++)
            {
                if (puzzles[p].Solve())
                {
                    puzzles[p].Print(Console.Out);
                    if (p < puzzles.Count - 1)
                        Console.Out.WriteLine();
                }
                else
                {
                    Console.Error.WriteLine($"puzzle {p + 1} has no solution");
                }
            }
        }

        static void Main(string[] args)
        {
            try
            {
                if (getBoolFlag(args, "help") || getBoolFlag(args, "h"))
                    Console.WriteLine(helpString);
                else if (getBoolFlag(args, "dl"))
                    runDownload(args);
                else
                    runSolve(args);
            }
            catch (Exception e)
            {
                Console.Error.WriteLine(e.Message);
            }
        }
    }
}
